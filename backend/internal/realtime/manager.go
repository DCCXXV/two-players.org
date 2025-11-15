package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/DCCXXV/twoplayers/backend/internal/games"
	appLogger "github.com/DCCXXV/twoplayers/backend/internal/logger"
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgconn"
)

type Manager struct {
	config            *config.Config
	connectionService service.ConnectionService
	roomService       service.RoomService
	playerService     service.PlayerService
	upgrader          websocket.Upgrader
	mu                sync.RWMutex
	clients           map[uuid.UUID]*Client
	rooms             map[uuid.UUID]*Room
	logger            *slog.Logger
}

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
		logger:  appLogger.Get(),
	}

	m.logger.Info("WebSocket manager initialized")
	m.StartCleanupTask(5 * time.Minute)
	m.CleanupStaleConnections()
	return m, nil
}

func (m *Manager) CleanupStaleConnections() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connections, err := m.connectionService.ListActiveConnections(ctx)
	if err != nil {
		m.logger.Error("Failed to list active connections for cleanup", "error", err)
		return
	}

	for _, conn := range connections {
		err := m.connectionService.DeleteConnection(ctx, conn.DisplayName)
		if err != nil {
			m.logger.Error("Failed to delete stale connection", "display_name", conn.DisplayName, "error", err)
		}
	}

	if len(connections) > 0 {
		m.logger.Info("Cleaned up stale connections", "count", len(connections))
	}
}

func (m *Manager) ServeWebSocket(w http.ResponseWriter, r *http.Request) error {
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.logger.Error("WebSocket upgrade failed", "error", err)
		return err
	}

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
		m.logger.Error("Failed to register connection", "error", dbErr)
		return fmt.Errorf("failed to register connection in DB: %w", dbErr)
	}
	if dbErr != nil {
		conn.Close()
		m.logger.Error("Failed to register connection after retries", "error", dbErr)
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

func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	m.clients[client.id] = client
	m.mu.Unlock()
	m.broadcastConnections()
}

func (m *Manager) unregisterClient(client *Client) {
	var roomToDelete *Room
	if client.currentRoom != nil {
		room := client.currentRoom
		if room.removeClient(client) {
			roomToDelete = room
		}
	}

	m.mu.Lock()
	if roomToDelete != nil {
		delete(m.rooms, roomToDelete.ID)
	}

	if _, ok := m.clients[client.id]; ok {
		delete(m.clients, client.id)
		close(client.send)

		if client.displayName != "" {
			go m.cleanupConnectionDB(client.displayName)
		}
	}
	m.mu.Unlock()
	m.broadcastConnections()
}

func (m *Manager) cleanupConnectionDB(displayName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.connectionService.DeleteConnection(ctx, displayName)
	if err != nil {
		m.logger.Error("Failed to delete connection from DB", "display_name", displayName, "error", err)
	}
}

func (m *Manager) StartCleanupTask(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			m.cleanupStaleRooms()
		}
	}()
}

func (m *Manager) cleanupStaleRooms() {
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
				delete(m.rooms, roomID)
			}
		}
		m.mu.Unlock()
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
	m.mu.RLock()
	defer m.mu.RUnlock()

	clientsByDisplayName := make(map[string]*Client)
	for _, client := range m.clients {
		clientsByDisplayName[client.displayName] = client
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connections, err := m.connectionService.ListActiveConnections(ctx)
	if err != nil {
		m.logger.Error("Failed to list active connections for broadcast", "error", err)
		return
	}

	payload := make([]map[string]interface{}, len(connections))

	for i, conn := range connections {
		status := "idle"
		var gameType *string

		if client, ok := clientsByDisplayName[conn.DisplayName]; ok {
			if client.currentRoom != nil {
				status = "in-game"
				gameType = &client.currentRoom.GameType
			}
		}

		payload[i] = map[string]interface{}{
			"display_name": conn.DisplayName,
			"status":       status,
			"game_type":    gameType,
		}
	}

	message := map[string]interface{}{
		"type":    "connections_update",
		"payload": payload,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		m.logger.Error("Failed to marshal connections update", "error", err)
		return
	}

	for _, client := range m.clients {
		client.send <- msgBytes
	}
}
