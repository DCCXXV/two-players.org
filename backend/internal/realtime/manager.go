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

	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/DCCXXV/twoplayers/backend/internal/games" // Importar el nuevo paquete de juegos
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
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
	currentRoom *Room
	send        chan []byte

	// Estado en la sala actual
	role     string // "player_0", "player_1", "spectator"
	joinedAt time.Time
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
	rooms   map[uuid.UUID]*Room
}

// GameInstance es una interfaz que cada juego debe implementar.
// Ahora se refiere a la interfaz games.Game del paquete games.
type GameInstance = games.Game

// Room representa una sala de juego activa en memoria.
type Room struct {
	ID         uuid.UUID
	GameType   string
	HostName   string
	Game       GameInstance
	manager    *Manager
	mu         sync.RWMutex
	Clients    map[uuid.UUID]*Client
	MaxPlayers int
}

// Métodos helper SIN locks para uso interno
func (r *Room) getPlayersInternal() []*Client {
	// NO usar locks aquí - se asume que el caller ya tiene el lock
	var players []*Client
	for _, client := range r.Clients {
		if strings.HasPrefix(client.role, "player_") {
			players = append(players, client)
		}
	}
	return players
}

func (r *Room) getSpectatorsInternal() []*Client {
	// NO usar locks aquí - se asume que el caller ya tiene el lock
	var spectators []*Client
	for _, client := range r.Clients {
		if client.role == "spectator" {
			spectators = append(spectators, client)
		}
	}
	return spectators
}

func (r *Room) getPlayerCountInternal() int {
	return len(r.getPlayersInternal())
}

func (r *Room) getSpectatorCountInternal() int {
	return len(r.getSpectatorsInternal())
}

// Métodos públicos CON locks (mantener para uso externo)
func (r *Room) GetPlayers() []*Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getPlayersInternal()
}

func (r *Room) GetSpectators() []*Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getSpectatorsInternal()
}

func (r *Room) GetPlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getPlayerCountInternal()
}

func (r *Room) GetSpectatorCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getSpectatorCountInternal()
}

// WebSocketMessage es la estructura genérica para la comunicación.
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
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

	// Lógica para crear una conexión y un nombre de usuario único
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

	// Creación del cliente y lanzamiento de los pumps
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
	log.Printf("Realtime State: Attempting to unregister client %s (%s). Current room: %v", client.id, client.displayName, client.currentRoom != nil)
	if client.currentRoom != nil {
		log.Printf("Realtime State: Client %s (%s) was in room %s. Calling removeClient.", client.displayName, client.id, client.currentRoom.ID)
		client.currentRoom.removeClient(client)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("Realtime State: Before unregistering client %s, total clients: %d", client.id, len(m.clients))
	if _, ok := m.clients[client.id]; ok {
		delete(m.clients, client.id)
		close(client.send)
		log.Printf("Realtime State: Successfully unregistered client %s (%s). Total clients now: %d", client.id, client.displayName, len(m.clients))

		if client.displayName != "" {
			go m.cleanupConnectionDB(client.displayName)
		}
	} else {
		log.Printf("Realtime State: Client %s (%s) not found in manager's clients map.", client.id, client.displayName)
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

func (m *Manager) getMaxPlayersForGame(gameType string) int {
	switch gameType {
	case "tic-tac-toe":
		return 2
	default:
		return 2
	}
}

// handleJoinRoom es el handler para el mensaje "join_room".
func (m *Manager) handleJoinRoom(client *Client, payload json.RawMessage) {
	log.Printf("🔄 handleJoinRoom: Starting for client %s", client.displayName)

	var req struct {
		RoomID string `json:"roomId"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Printf("❌ handleJoinRoom: Invalid payload: %v", err)
		client.sendError("Invalid payload for join_room")
		return
	}

	log.Printf("🔄 handleJoinRoom: Parsed roomId: %s", req.RoomID)

	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		log.Printf("❌ handleJoinRoom: Invalid UUID: %v", err)
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
		// La sala no está en memoria, la cargamos desde la BD
		dbRoom, err := m.roomService.GetRoomByID(context.Background(), roomID)
		if err != nil {
			m.mu.Unlock()
			client.sendError("Room not found in database.")
			return
		}

		var game GameInstance
		game, err = games.NewGame(dbRoom.GameType)
		if err != nil {
			m.mu.Unlock()
			client.sendError(fmt.Sprintf("Error creating game instance: %v", err))
			return
		}

		// Crear la sala en memoria
		room = &Room{
			ID:         uuid.UUID(dbRoom.ID.Bytes),
			GameType:   dbRoom.GameType,
			Clients:    make(map[uuid.UUID]*Client),
			HostName:   dbRoom.HostDisplayName,
			Game:       game,
			manager:    m,
			MaxPlayers: m.getMaxPlayersForGame(dbRoom.GameType),
		}
		m.rooms[room.ID] = room
		log.Printf("Activated room %s in memory.", room.ID)
	}
	m.mu.Unlock()

	log.Printf("🔄 handleJoinRoom: About to call room.addClient")
	room.addClient(client)
	log.Printf("✅ handleJoinRoom: Completed successfully")
}

// -----------------------------------------------------
// MÉTODOS DE LA SALA (ROOM)
// -----------------------------------------------------

func (r *Room) addClient(client *Client) {
	log.Printf("🔄 addClient: Starting for client %s in room %s", client.displayName, r.ID)

	r.mu.Lock()

	if _, ok := r.Clients[client.id]; ok {
		log.Printf("⚠️  addClient: Client already in room")
		r.mu.Unlock()
		return
	}

	// Determinar rol basado en orden de llegada - USAR MÉTODOS INTERNOS
	playerCount := r.getPlayerCountInternal() // ← Sin deadlock
	log.Printf("🔄 addClient: Current player count: %d", playerCount)

	var playerOrder int16
	isPlayer := false
	if playerCount == 0 {
		client.role = "player_0"
		playerOrder = 0
		isPlayer = true
	} else if playerCount == 1 {
		client.role = "player_1"
		playerOrder = 1
		isPlayer = true
	} else {
		client.role = "spectator"
	}

	// Si el cliente es un jugador, lo persistimos en la BD antes de añadirlo a la sala en memoria.
	if isPlayer {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		params := service.CreatePlayerParams{
			RoomID:            pgtype.UUID{Bytes: r.ID, Valid: true},
			PlayerDisplayName: client.displayName,
			PlayerOrder:       playerOrder,
		}
		_, err := r.manager.playerService.CreatePlayer(ctx, params)
		if err != nil {
			var pgErr *pgconn.PgError
			// Si es una violación de restricción única, no es un error fatal aquí.
			// Significa que el jugador ya estaba registrado (p. ej., el anfitrión que vuelve a unirse).
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				log.Printf("ℹ️ addClient: Player %s already exists in room %s. Continuing.", client.displayName, r.ID)
			} else {
				log.Printf("❌ addClient: Failed to create player in DB for client %s in room %s: %v", client.displayName, r.ID, err)
				client.sendError("Failed to join room: could not save player data.")
				r.mu.Unlock()
				return
			}
		} else {
			log.Printf("✅ addClient: Successfully created player record in DB for %s as player %d", client.displayName, playerOrder)
		}
	}

	client.joinedAt = time.Now()
	client.currentRoom = r
	r.Clients[client.id] = client

	log.Printf("✅ addClient: Client %s joined room %s as %s (total players: %d)",
		client.displayName, r.ID, client.role, r.getPlayerCountInternal())

	client.sendMessage("join_success", map[string]any{
		"roomId": r.ID.String(),
		"role":   client.role,
	})

	r.mu.Unlock() // ← Liberar antes del broadcast

	log.Printf("🔄 addClient: About to call broadcastRoomState")
	r.broadcastRoomState()
	log.Printf("✅ addClient: Completed successfully")
}

func (r *Room) removeClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock() // Ensure unlock happens

	log.Printf("Room %s: removeClient called for client %s (%s). Host: %s. Current clients in room: %d", r.ID, client.displayName, client.id, r.HostName, len(r.Clients))

	if _, ok := r.Clients[client.id]; !ok {
		log.Printf("Room %s: Client %s not found in room. Aborting removeClient.", r.ID, client.displayName)
		return // Client not in room
	}

	isHost := client.displayName == r.HostName
	log.Printf("Room %s: Client %s is host: %t", r.ID, client.displayName, isHost)

	delete(r.Clients, client.id)
	client.currentRoom = nil
	log.Printf("Room %s: Client %s removed. Remaining clients in room: %d", r.ID, client.displayName, len(r.Clients))

	// If the host leaves, the room is closed for everyone.
	if isHost {
		log.Printf("Room %s: Host %s left. Closing room for all remaining clients.", r.ID, client.displayName)

		// Notify remaining clients and clear their room association
		for _, otherClient := range r.Clients {
			log.Printf("Room %s: Notifying client %s about room closure.", r.ID, otherClient.displayName)
			otherClient.sendMessage("room_closed", map[string]string{"message": "The host has left the room."})
			otherClient.currentRoom = nil
		}

		// Clear the clients map for this room
		r.Clients = make(map[uuid.UUID]*Client) // This effectively empties the map

		log.Printf("Room %s: All clients notified and room client map cleared. Remaining clients in room (after clear): %d", r.ID, len(r.Clients))

		// Remove the room from the manager's list of active rooms
		r.manager.mu.Lock()
		log.Printf("Manager: Before deleting room %s, total rooms: %d", r.ID, len(r.manager.rooms))
		delete(r.manager.rooms, r.ID)
		log.Printf("Manager: Room %s deleted. Total rooms now: %d", r.ID, len(r.manager.rooms))
		r.manager.mu.Unlock()
		log.Printf("Room %s closed and removed from memory because host left.", r.ID)

	} else if len(r.Clients) == 0 {
		log.Printf("Room %s: Last non-host client %s left. Room is now empty.", r.ID, client.displayName)
		// If the last client (who wasn't the host) leaves, the room is also removed.
		r.manager.mu.Lock()
		log.Printf("Manager: Before deleting room %s, total rooms: %d", r.ID, len(r.manager.rooms))
		delete(r.manager.rooms, r.ID)
		log.Printf("Manager: Room %s deleted. Total rooms now: %d", r.ID, len(r.manager.rooms))
		r.manager.mu.Unlock()
		log.Printf("Room %s is now empty and has been removed from memory.", r.ID)
	} else {
		// A non-host client left, just update the state for everyone else.
		log.Printf("Room %s: Non-host client %s left. Broadcasting updated room state.", r.ID, client.displayName)
		r.broadcastRoomState()
	}
}

func (r *Room) broadcastRoomState() {
	log.Printf("🔄 broadcastRoomState: Starting for room %s", r.ID)

	r.mu.RLock()
	players := r.getPlayersInternal()       // ← Usar métodos internos
	spectators := r.getSpectatorsInternal() // ← Usar métodos internos

	// Convertir a nombres para el frontend
	playerNames := make([]string, len(players))
	for i, p := range players {
		playerNames[i] = p.displayName
	}

	spectatorNames := make([]string, len(spectators))
	for i, s := range spectators {
		spectatorNames[i] = s.displayName
	}

	gameState := r.Game.GetGameState()
	r.mu.RUnlock() // ← Liberar antes de crear el mapa

	// Estado combinado
	roomState := map[string]any{
		"roomId":         r.ID.String(),
		"gameType":       r.GameType,
		"players":        playerNames,
		"spectators":     spectatorNames,
		"playerCount":    len(players),
		"spectatorCount": len(spectators),
		"maxPlayers":     r.MaxPlayers,
		"canStart":       len(players) == r.MaxPlayers,
		"game":           gameState,
	}

	log.Printf("🔄 broadcastRoomState: Room state created: %+v", roomState)

	r.broadcastMessage("game_state_update", roomState)

	log.Printf("✅ broadcastRoomState: Completed")
}

func (r *Room) broadcastMessage(msgType string, payload any) {
	log.Printf("🔄 broadcastMessage: Creating message type '%s' for room %s", msgType, r.ID)

	message, err := createWebSocketMessage(msgType, payload)
	if err != nil {
		log.Printf("❌ broadcastMessage: Failed to create message '%s': %v", msgType, err)
		return
	}

	r.mu.RLock()
	clientsCopy := make([]*Client, 0, len(r.Clients))
	for _, client := range r.Clients {
		clientsCopy = append(clientsCopy, client)
	}
	r.mu.RUnlock()

	log.Printf("🔄 broadcastMessage: Broadcasting '%s' to %d clients", msgType, len(clientsCopy))

	for i, client := range clientsCopy {
		select {
		case client.send <- message:
			log.Printf("✅ broadcastMessage: Sent '%s' to client %d: %s", msgType, i, client.displayName)
		default:
			log.Printf("⚠️  broadcastMessage: Client %s send channel full for message '%s'", client.displayName, msgType)
		}
	}

	log.Printf("✅ broadcastMessage: Completed broadcasting '%s'", msgType)
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
		case "make_move":
			// Manejar movimientos del juego
			if c.currentRoom != nil {
				c.handleGameMove(msg.Payload)
			} else {
				c.sendError("Not in a room.")
			}
		default:
			c.sendError(fmt.Sprintf("Unknown message type '%s'.", msg.Type))
		}
	}
}

func (c *Client) handleGameMove(payload json.RawMessage) {
	// Determinar índice del jugador
	var playerIndex int
	switch c.role {
	case "player_0":
		playerIndex = 0
	case "player_1":
		playerIndex = 1
	default:
		c.sendError("Spectators cannot make moves.")
		return
	}

	// Parsear el movimiento (esto dependerá del tipo de juego)
	var move any
	if err := json.Unmarshal(payload, &move); err != nil {
		c.sendError("Invalid move format.")
		return
	}

	// Intentar hacer el movimiento
	if err := c.currentRoom.Game.HandleMove(playerIndex, move); err != nil {
		c.sendError(err.Error())
		return
	}

	// Broadcast del nuevo estado
	c.currentRoom.broadcastRoomState()
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


