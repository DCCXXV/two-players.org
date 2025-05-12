package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn        *websocket.Conn
	manager     *Manager
	id          uuid.UUID
	displayName string
	// currentRoomID uuid.UUID
	send chan []byte
}

type Manager struct {
	config            *config.Config
	connectionService service.ConnectionService

	upgrader websocket.Upgrader

	mu           sync.RWMutex
	clients      map[uuid.UUID]*Client
	nameToClient map[string]*Client
	rooms        map[uuid.UUID]*Room
}

type GameInstance interface {
	HandleMessage(clientID uuid.UUID, messageType string, payload json.RawMessage) error
	AddPlayer(client *Client) error
	RemovePlayer(clientID uuid.UUID) error
	GetStateForPlayer(clientID uuid.UUID) interface{}
}

type Room struct {
	ID       uuid.UUID
	Name     string
	GameType string
	Clients  map[uuid.UUID]*Client
	HostID   uuid.UUID
	Game     GameInstance
	manager  *Manager
	mu       sync.RWMutex
}

func NewManager(cfg *config.Config, cs service.ConnectionService, rm service.RoomService, pl service.PlayerService) (*Manager, error) {
	m := &Manager{
		config:            cfg,
		connectionService: cs,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					log.Printf("WebSocket CheckOrigin: Allowing request with empty Origin header.")
					return true
				}
				if origin == cfg.AllowedOrigins {
					log.Printf("WebSocket CheckOrigin: Allowing origin: %s", origin)
					return true
				}
				log.Printf("WARN: WebSocket CheckOrigin: Denied origin: %s", origin)
				return false
			},
		},
		clients: make(map[uuid.UUID]*Client),
	}

	log.Println("Realtime Manager (WebSocket) initialized")
	return m, nil
}

func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client.id] = client
	log.Printf("Realtime State: Registered client %s", client.id)
}

func (m *Manager) unregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[client.id]; ok {
		delete(m.clients, client.id)
		close(client.send)
		log.Printf("Realtime State: Unregistered client %s", client.id)

		if client.displayName != "" {
			go m.cleanupConnectionDB(client.displayName)
		} else {
			log.Printf("WARN: Client %s disconnected before having a display name.", client.id)
		}
	}
}

func (m *Manager) cleanupConnectionDB(displayName string) {
	log.Printf("Realtime Cleanup: Attempting to delete connection for %s from DB", displayName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.connectionService.DeleteConnection(ctx, displayName)
	if err != nil {
		log.Printf("ERROR: Realtime Cleanup: Failed to delete connection for %s from DB: %v", displayName, err)
	} else {
		log.Printf("Realtime Cleanup: Deleted connection for %s from DB.", displayName)
	}
}

func (m *Manager) sendToClient(clientID uuid.UUID, message []byte) {
	m.mu.RLock()
	client, ok := m.clients[clientID]
	m.mu.RUnlock()

	if ok {
		select {
		case client.send <- message:
		default:
			log.Printf("WARN: Send channel full for client %s. Message dropped.", clientID)
			client.conn.Close()
			m.unregisterClient(client)
		}
	} else {
		log.Printf("WARN: Attempted to send message to non-existent client ID: %s", clientID)
	}
}

func (m *Manager) broadcast(message []byte, excludeClientID *uuid.UUID) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for id, client := range m.clients {
		if excludeClientID != nil && id == *excludeClientID {
			continue
		}
		select {
		case client.send <- message:
		default:
			log.Printf("WARN: Send channel full during broadcast for client %s. Message dropped.", id)
		}
	}
}

// TODO: Add methods to manage joinRoom, leaveRoom, broadcastToRoom

// Placeholder
func (m *Manager) RunCleanupTask() {
	log.Println("Realtime: Cleanup task placeholder started.")
}

var namePrefixes = []string{"Alice", "Bob"}

const maxGenerateNameRetries = 5

func generateAliceOrBobName() string {
	prefixIndex := rand.Intn(len(namePrefixes))
	prefix := namePrefixes[prefixIndex]
	randomNumber := rand.Intn(9000) + 1000
	return fmt.Sprintf("%s#%d", prefix, randomNumber)
}

func (m *Manager) ServeWebSocket(w http.ResponseWriter, r *http.Request) error {
	log.Println("Realtime: Received WebSocket upgrade request")

	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: WebSocket upgrade failed: %v", err)
		return err
	}
	log.Printf("Realtime: WebSocket upgrade successful for %s", conn.RemoteAddr())

	var dbErr error
	var generatedName string

	for i := 0; i < maxGenerateNameRetries; i++ {
		generatedName = generateAliceOrBobName()
		log.Printf("Realtime: Attempting generated name %s for conn %s (try %d)", generatedName, conn.RemoteAddr(), i+1)

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		_, dbErr = m.connectionService.CreateConnection(ctx, service.CreateConnectionParams{
			DisplayName: generatedName,
		})

		cancel()

		if dbErr == nil {
			log.Printf("Realtime: Successfully registered connection %s with name %s in DB", conn.RemoteAddr(), generatedName)
			break
		}

		var pgErr *pgconn.PgError
		if errors.As(dbErr, &pgErr) && pgErr.Code == "23505" {
			log.Printf("Realtime: Generated name %s collision, retrying...", generatedName)
			continue
		} else {
			log.Printf("ERROR: Realtime: Failed to create connection in DB for %s after generating name %s: %v", conn.RemoteAddr(), generatedName, dbErr)
			conn.Close()
			return fmt.Errorf("failed to register connection in DB: %w", dbErr)
		}
	}

	if dbErr != nil {
		log.Printf("ERROR: Realtime: Failed to register connection %s after %d retries: %v", conn.RemoteAddr(), maxGenerateNameRetries, dbErr)
		conn.Close()
		return fmt.Errorf("failed to register connection after retries: %w", dbErr)
	}

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

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

func (c *Client) sendConnectionReady() {
	msg := WebSocketMessage{
		Type:    "connection_ready",
		Payload: map[string]string{"displayName": c.displayName},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ERROR: Failed to marshal connection_ready message for client %s: %v", c.id, err)
		return
	}

	select {
	case c.send <- jsonData:
		log.Printf("Realtime: Sent connection_ready to client %s (%s)", c.id, c.displayName)
	default:
		log.Printf("WARN: Send channel full for client %s on connection_ready. Message dropped.", c.id)
	}
}

func (c *Client) readPump() {

	defer func() {
		log.Printf("Realtime: Closing readPump for client %s (%s)", c.id, c.displayName)
		c.manager.unregisterClient(c)
		c.conn.Close()
	}()

	// c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR: WebSocket read error for client %s (%s): %v", c.id, c.displayName, err)
			} else {
				log.Printf("INFO: WebSocket closed for client %s (%s): %v", c.id, c.displayName, err)
			}
			break
		}

		if messageType != websocket.TextMessage {
			log.Printf("WARN: Received non-text message type %d from client %s", messageType, c.id)
			continue
		}

		log.Printf("Realtime: Received message from %s (%s): %s", c.id, c.displayName, string(message))

		// TODO: Parsear el mensaje (ej: JSON) y actuar según su tipo/contenido
		// var msg WebSocketMessage
		// if err := json.Unmarshal(message, &msg); err != nil {
		//     log.Printf("ERROR: Failed to unmarshal message from client %s: %v", c.id, err)
		//     continue
		// }
		// switch msg.Type {
		// case "change_name":
		//     // Llama a m.handleChangeName(c, msg.Payload)
		// case "create_room":
		//     // Llama a m.handleCreateRoom(c, msg.Payload)
		// // ... otros casos ...
		// default:
		//     log.Printf("WARN: Unknown message type '%s' from client %s", msg.Type, c.id)
		// }
	}
}

func (c *Client) writePump() {

	defer func() {
		log.Printf("Realtime: Closing writePump for client %s (%s)", c.id, c.displayName)
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Printf("Realtime: Send channel closed for client %s (%s). Closing connection.", c.id, c.displayName)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("ERROR: WebSocket write error for client %s (%s): %v", c.id, c.displayName, err)
				return
			}
			log.Printf("Realtime: Sent message to client %s (%s): %s", c.id, c.displayName, string(message))
		}
	}
}
