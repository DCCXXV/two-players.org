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
	"github.com/DCCXXV/twoplayers/backend/internal/games"
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgconn"
)

// Manager is the orchestrator of all real-time connections and rooms.
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

// GameInstance is an interface that each game must implement.
// It now refers to the games.Game interface from the games package.
type GameInstance = games.Game

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
	m.StartCleanupTask(5 * time.Minute) // Start periodic cleanup task
	m.CleanupStaleConnections()
	return m, nil
}

func (m *Manager) CleanupStaleConnections() {
	log.Println("Cleaning up stale connections from the database...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connections, err := m.connectionService.ListActiveConnections(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to list active connections for cleanup: %v", err)
		return
	}

	if len(connections) == 0 {
		log.Println("No stale connections to clean up.")
		return
	}

	for _, conn := range connections {
		err := m.connectionService.DeleteConnection(ctx, conn.DisplayName)
		if err != nil {
			log.Printf("ERROR: Failed to delete stale connection for %s: %v", conn.DisplayName, err)
		} else {
			log.Printf("Successfully cleaned up stale connection for %s.", conn.DisplayName)
		}
	}
	log.Printf("Finished cleaning up %d stale connections.", len(connections))
}

// ServeWebSocket is the entry point for a new connection.
func (m *Manager) ServeWebSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ERROR: WebSocket upgrade failed: %v", err)
		return err
	}

	// Logic to create a connection and a unique username
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

	// Create client and launch pumps
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
	m.clients[client.id] = client
	m.mu.Unlock()
	log.Printf("Realtime State: Registered client %s (%s)", client.id, client.displayName)
	m.broadcastConnections()
}

func (m *Manager) unregisterClient(client *Client) {
	log.Printf("Realtime State: Attempting to unregister client %s (%s). Current room: %v", client.id, client.displayName, client.currentRoom != nil)

	var roomToDelete *Room
	if client.currentRoom != nil {
		room := client.currentRoom
		// removeClient will only lock the room, not the manager, and will tell us if it needs to be deleted.
		if room.removeClient(client) {
			roomToDelete = room
		}
	}

	m.mu.Lock()
	// If a room was marked for deletion, remove it from the manager.
	if roomToDelete != nil {
		log.Printf("Manager: Deleting room %s because it's empty or the host left.", roomToDelete.ID)
		delete(m.rooms, roomToDelete.ID)
		log.Printf("Manager: Room %s deleted. Total rooms now: %d", roomToDelete.ID, len(m.rooms))
	}

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
	m.mu.Unlock()
	m.broadcastConnections()
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

// StartCleanupTask starts a goroutine to periodically clean up empty rooms.
func (m *Manager) StartCleanupTask(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			m.cleanupStaleRooms()
		}
	}()
	log.Printf("Periodic cleanup task started. Interval: %s", interval)
}

// cleanupStaleRooms iterates over the rooms and removes those that are empty to prevent memory leaks.
func (m *Manager) cleanupStaleRooms() {
	log.Println("Realtime Cleanup: Running periodic task for stale rooms...")
	var roomsToDelete []uuid.UUID

	m.mu.RLock()
	for roomID, room := range m.rooms {
		room.mu.RLock()
		clientCount := len(room.Clients)
		room.mu.RUnlock()

		if clientCount == 0 {
			roomsToDelete = append(roomsToDelete, roomID)
		}
	}
	m.mu.RUnlock()

	if len(roomsToDelete) > 0 {
		m.mu.Lock()
		for _, roomID := range roomsToDelete {
			if _, ok := m.rooms[roomID]; ok {
				log.Printf("Realtime Cleanup: Removing empty/stale room %s", roomID)
				delete(m.rooms, roomID)
			}
		}
		m.mu.Unlock()
		log.Printf("Realtime Cleanup: Finished removing %d stale rooms.", len(roomsToDelete))
	} else {
		log.Println("Realtime Cleanup: No stale rooms found to remove.")
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

func generateAliceOrBobName() string {
	namePrefixes := []string{"Alice", "Bob"}
	prefix := namePrefixes[rand.Intn(len(namePrefixes))]
	return fmt.Sprintf("%s#%d", prefix, rand.Intn(9000)+1000)
}

func (m *Manager) broadcastConnections() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connections, err := m.connectionService.ListActiveConnections(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to list active connections for broadcast: %v", err)
		return
	}

	message := map[string]interface{}{
		"type":    "connections_update",
		"payload": connections,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("ERROR: Failed to marshal connections update: %v", err)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, client := range m.clients {
		client.send <- msgBytes
	}
}
