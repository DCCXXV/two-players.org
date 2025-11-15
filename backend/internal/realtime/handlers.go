package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/games"
	"github.com/google/uuid"
)

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

		room = &Room{
			ID:              uuid.UUID(dbRoom.ID.Bytes),
			GameType:        dbRoom.GameType,
			Clients:         make(map[uuid.UUID]*Client),
			HostName:        dbRoom.HostDisplayName,
			Game:            game,
			manager:         m,
			MaxPlayers:      m.getMaxPlayersForGame(dbRoom.GameType),
			rematchRequests: make(map[uuid.UUID]bool),
		}
		m.rooms[room.ID] = room
	}
	m.mu.Unlock()

	room.addClient(client)
}

func (c *Client) handleGameMove(payload json.RawMessage) {
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

	var move any
	if err := json.Unmarshal(payload, &move); err != nil {
		c.sendError("Invalid move format.")
		return
	}

	if err := c.currentRoom.Game.HandleMove(playerIndex, move); err != nil {
		c.sendError(err.Error())
		return
	}

	c.currentRoom.broadcastRoomState()
}

func (r *Room) handleRematch(client *Client) {
	r.mu.Lock()

	if !r.Game.IsGameOver() {
		client.sendError("The game is not over yet.")
		r.mu.Unlock()
		return
	}

	r.rematchRequests[client.id] = true

	players := r.getPlayersInternal()
	rematchCount := len(r.rematchRequests)
	allPlayersRequestedRematch := rematchCount == len(players) && len(players) == r.MaxPlayers

	r.mu.Unlock()

	r.broadcastRoomState()

	if allPlayersRequestedRematch {
		r.mu.Lock()
		r.Game.Reset()

		var player0, player1 *Client
		for _, p := range r.getPlayersInternal() {
			if p.role == "player_0" {
				player0 = p
			} else if p.role == "player_1" {
				player1 = p
			}
		}
		if player0 != nil && player1 != nil {
			player0.role, player1.role = player1.role, player0.role
		}

		r.rematchRequests = make(map[uuid.UUID]bool)
		r.mu.Unlock()

		go r.broadcastRoomState()
	}
}

func (r *Room) handleChatMessage(client *Client, payload json.RawMessage) {
	var req struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		client.sendError("Invalid chat message payload.")
		return
	}

	if req.Message == "" {
		return
	}

	chatMessage := map[string]any{
		"displayName": client.displayName,
		"message":     req.Message,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	r.broadcastMessage("chat_message", chatMessage)
}

func (m *Manager) handleUpdateDisplayName(client *Client, payload json.RawMessage) {
	var req struct {
		DisplayName string `json:"displayName"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		client.sendError("Invalid payload for update_display_name")
		return
	}

	newName := req.DisplayName
	if newName == "" {
		client.sendError("Display name cannot be empty")
		return
	}

	oldName := client.displayName
	client.displayName = newName

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.connectionService.UpdateConnectionName(ctx, oldName, newName)
	if err != nil {
		m.logger.Error("Failed to update display name in DB", "old_name", oldName, "new_name", newName, "error", err)
		client.displayName = oldName
		client.sendError("Failed to update display name")
		return
	}

	m.logger.Info("Display name updated", "old_name", oldName, "new_name", newName, "client_id", client.id)
	client.sendMessage("connection_ready", map[string]string{"displayName": newName})
}
