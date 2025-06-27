package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/config" // Necesario para los structs de la BD
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgconn"
)

// -----------------------------------------------------
// ESTRUCTURAS PRINCIPALES
// -----------------------------------------------------

// Client representa un único usuario conectado vía WebSocket.
type Client struct {
	conn        *websocket.Conn
	manager     *Manager
	id          uuid.UUID
	displayName string
	currentRoom *Room // Referencia a la sala en la que se encuentra el cliente.
	send        chan []byte
}

// Manager es el director de orquesta de todas las conexiones y salas en tiempo real.
type Manager struct {
	config            *config.Config
	connectionService service.ConnectionService
	roomService       service.RoomService
	playerService     service.PlayerService

	upgrader websocket.Upgrader

	mu      sync.RWMutex
	clients map[uuid.UUID]*Client
	rooms   map[uuid.UUID]*Room // Mapa de salas activas en memoria.
}

// GameInstance es una interfaz que cada juego (Tres en Raya, etc.) debe implementar.
type GameInstance interface {
	HandleMessage(client *Client, payload json.RawMessage) error
	AddPlayer(client *Client) error
	RemovePlayer(client *Client)
	GetState() any
}

// Room representa una sala de juego activa en memoria.
type Room struct {
	ID       uuid.UUID
	GameType string
	Clients  map[uuid.UUID]*Client
	HostName string
	Game     GameInstance
	manager  *Manager
	mu       sync.RWMutex
}

// WebSocketMessage es la estructura genérica para la comunicación.
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"` // Usamos RawMessage para retrasar el parseo del payload.
}

// ErrorPayload es una estructura estándar para enviar errores al cliente.
type ErrorPayload struct {
	Message string `json:"message"`
}

// -----------------------------------------------------
// CONSTRUCTOR Y MÉTODOS DEL MANAGER
// -----------------------------------------------------

func NewManager(cfg *config.Config, cs service.ConnectionService, rs service.RoomService, ps service.PlayerService) (*Manager, error) {
	allowedOriginsSlice := strings.Split(cfg.AllowedOrigins, ",")
	m := &Manager{
		config:            cfg,
		connectionService: cs,
		roomService:       rs,
		playerService:     ps,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				for _, allowed := range allowedOriginsSlice {
					if origin == allowed {
						return true
					}
				}
				return false
			},
		},
		clients: make(map[uuid.UUID]*Client),
		rooms:   make(map[uuid.UUID]*Room),
	}

	log.Println("Realtime Manager (WebSocket) initialized")
	return m, nil
}

// ServeWebSocket es el punto de entrada para una nueva conexión.
func (m *Manager) ServeWebSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: WebSocket upgrade failed: %v", err)
		return err
	}

	// --- Lógica para crear una conexión y un nombre de usuario único ---
	var dbErr error
	var generatedName string
	for i := 0; i < 5; i++ {
		generatedName = generateAliceOrBobName()
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		_, dbErr = m.connectionService.CreateConnection(ctx, service.CreateConnectionParams{
			DisplayName: generatedName,
		})
		cancel()
		if dbErr == nil {
			break
		}
		var pgErr *pgconn.PgError
		if errors.As(dbErr, &pgErr) && pgErr.Code == "23505" {
			continue
		}
		conn.Close()
		return fmt.Errorf("failed to register connection in DB: %w", dbErr)
	}
	if dbErr != nil {
		conn.Close()
		return fmt.Errorf("failed to register connection after retries: %w", dbErr)
	}

	// --- Creación del cliente y lanzamiento de los pumps ---
	client := &Client{
		conn:        conn,
		manager:     m,
		id:          uuid.New(),
		displayName: generatedName,
		send:        make(chan []byte, 256),
	}

	m.registerClient(client)
	client.sendConnectionReady()

	go client.writePump()
	go client.readPump()

	return nil
}

func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client.id] = client
	log.Printf("Realtime State: Registered client %s (%s)", client.id, client.displayName)
}

func (m *Manager) unregisterClient(client *Client) {
	if client.currentRoom != nil {
		client.currentRoom.removeClient(client)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client.id]; ok {
		delete(m.clients, client.id)
		close(client.send)
		log.Printf("Realtime State: Unregistered client %s (%s)", client.id, client.displayName)

		if client.displayName != "" {
			go m.cleanupConnectionDB(client.displayName)
		}
	}
}

func (m *Manager) cleanupConnectionDB(displayName string) {
	log.Printf("Realtime Cleanup: Deleting connection for %s from DB. This will CASCADE-delete their rooms.", displayName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.connectionService.DeleteConnection(ctx, displayName)
	if err != nil {
		log.Printf("ERROR: Realtime Cleanup: Failed to delete connection for %s from DB: %v", displayName, err)
	} else {
		log.Printf("Realtime Cleanup: Successfully deleted connection and associated data for %s from DB.", displayName)
	}
}

// handleJoinRoom es el handler para el mensaje "join_room". Activa una sala si es necesario.
func (m *Manager) handleJoinRoom(client *Client, payload json.RawMessage) {
	var req struct {
		RoomID string `json:"roomId"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		client.sendError("Invalid payload for join_room")
		return
	}
	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		client.sendError("Invalid room ID format")
		return
	}

	if client.currentRoom != nil && client.currentRoom.ID != roomID {
		client.sendError("You are already in another room.")
		return
	}

	m.mu.Lock()
	room, ok := m.rooms[roomID]
	if !ok {
		// La sala no está en memoria, la cargamos desde la BD.
		dbRoom, err := m.roomService.GetRoomByID(context.Background(), roomID)
		if err != nil {
			m.mu.Unlock()
			client.sendError("Room not found in database.")
			return
		}

		var game GameInstance
		switch dbRoom.GameType {
		case "tic-tac-toe":
			//game = games.NewTicTacToe() // ¡Asegúrate de que este constructor exista!
		default:
			m.mu.Unlock()
			client.sendError(fmt.Sprintf("Game type '%s' not supported", dbRoom.GameType))
			return
		}

		// Crear la sala en memoria, convirtiendo los tipos de la BD a tipos de Go.
		room = &Room{
			ID:       uuid.UUID(dbRoom.ID.Bytes),
			GameType: dbRoom.GameType,
			Clients:  make(map[uuid.UUID]*Client),
			HostName: dbRoom.HostDisplayName,
			Game:     game,
			manager:  m,
		}
		m.rooms[room.ID] = room
		log.Printf("Activated room %s in memory.", room.ID)
	}
	m.mu.Unlock()

	room.addClient(client)
}

// -----------------------------------------------------
// MÉTODOS DE LA SALA (ROOM)
// -----------------------------------------------------

func (r *Room) addClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Clients[client.id]; ok {
		return // El cliente ya está en la sala.
	}

	// Usar el servicio para registrar al jugador en la BD.
	joinResult, err := r.manager.roomService.JoinRoom(context.Background(), service.JoinRoomInput{
		RoomID:      r.ID,
		DisplayName: client.displayName,
	})
	if err != nil {
		log.Printf("ERROR: Failed to add player to room in DB: %v", err)
		client.sendError(fmt.Sprintf("Cannot join game: %s", err.Error()))
		return
	}

	if err := r.Game.AddPlayer(client); err != nil {
		client.sendError(fmt.Sprintf("Cannot join game logic: %s", err.Error()))
		return
	}

	r.Clients[client.id] = client
	client.currentRoom = r
	log.Printf("Client %s joined room %s as %s", client.displayName, r.ID, joinResult.Role)

	client.sendMessage("join_success", map[string]any{"roomId": r.ID.String()})
	r.broadcastStateUpdate()
}

func (r *Room) removeClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Clients[client.id]; ok {
		delete(r.Clients, client.id)
		client.currentRoom = nil
		r.Game.RemovePlayer(client)

		log.Printf("Client %s left room %s", client.displayName, r.ID)

		if len(r.Clients) == 0 {
			r.manager.mu.Lock()
			delete(r.manager.rooms, r.ID)
			r.manager.mu.Unlock()
			log.Printf("Room %s is now empty and has been removed.", r.ID)
		} else {
			r.broadcastStateUpdate()
		}
	}
}

func (r *Room) broadcastStateUpdate() {
	state := r.Game.GetState()
	message, err := createWebSocketMessage("game_state_update", state)
	if err != nil {
		log.Printf("ERROR: Failed to create game_state_update message: %v", err)
		return
	}

	// Itera sobre una copia del mapa de clientes para evitar problemas de concurrencia
	// si un cliente se desconecta mientras se está transmitiendo.
	r.mu.RLock()
	clientsCopy := make([]*Client, 0, len(r.Clients))
	for _, client := range r.Clients {
		clientsCopy = append(clientsCopy, client)
	}
	r.mu.RUnlock()

	for _, client := range clientsCopy {
		// Envío no bloqueante para evitar que un cliente lento bloquee a todos los demás.
		select {
		case client.send <- message:
		default:
			log.Printf("WARN: Client %s send channel is full. Dropping message.", client.id)
		}
	}
}

// -----------------------------------------------------
// MÉTODOS DEL CLIENTE (PUMPS Y HELPERS)
// -----------------------------------------------------

func (c *Client) readPump() {
	defer func() {
		log.Printf("Realtime: Closing readPump for client %s (%s)", c.id, c.displayName)
		c.manager.unregisterClient(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR: WebSocket read error for client %s: %v", c.id, err)
			}
			break
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			c.sendError("Invalid message format.")
			continue
		}

		switch msg.Type {
		case "join_room":
			c.manager.handleJoinRoom(c, msg.Payload)
		default:
			// Evita mantener el bloqueo de la sala durante la transmisión.
			// Primero, maneja el mensaje y luego, si es necesario, transmite el estado.
			if c.currentRoom != nil {
				var shouldBroadcast bool
				c.currentRoom.mu.Lock()
				err := c.currentRoom.Game.HandleMessage(c, msg.Payload)
				if err != nil {
					c.sendError(err.Error())
				} else {
					shouldBroadcast = true
				}
				c.currentRoom.mu.Unlock()

				if shouldBroadcast {
					c.currentRoom.broadcastStateUpdate()
				}
			} else {
				c.sendError(fmt.Sprintf("Unknown message type '%s' or not in a room.", msg.Type))
			}
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *Client) sendConnectionReady() {
	c.sendMessage("connection_ready", map[string]string{"displayName": c.displayName})
}

func (c *Client) sendMessage(msgType string, payload any) {
	jsonData, err := createWebSocketMessage(msgType, payload)
	if err != nil {
		log.Printf("ERROR: Failed to marshal message '%s' for client %s: %v", msgType, c.id, err)
		return
	}
	c.send <- jsonData
}

func (c *Client) sendError(message string) {
	c.sendMessage("error", ErrorPayload{Message: message})
}

// -----------------------------------------------------
// FUNCIONES DE AYUDA
// -----------------------------------------------------

func createWebSocketMessage(msgType string, payload any) ([]byte, error) {
	// Serializa el payload a json.RawMessage antes de construir el mensaje final.
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	msg := WebSocketMessage{Type: msgType, Payload: payloadBytes}
	return json.Marshal(msg)
}

func generateAliceOrBobName() string {
	namePrefixes := []string{"Alice", "Bob"}
	prefix := namePrefixes[rand.Intn(len(namePrefixes))]
	return fmt.Sprintf("%s#%d", prefix, rand.Intn(9000)+1000)
}
