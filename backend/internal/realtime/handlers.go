package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/DCCXXV/twoplayers/backend/internal/games"
	"github.com/google/uuid"
)

// handleJoinRoom is the handler for the "join_room" message.
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
		// The room is not in memory, load it from the DB
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

		// Create the room in memory
		room = &Room{
			ID:              uuid.UUID(dbRoom.ID.Bytes),
			GameType:        dbRoom.GameType,
			Clients:         make(map[uuid.UUID]*Client),
			HostName:        dbRoom.HostDisplayName,
			Game:            game,
			manager:         m,
			MaxPlayers:      m.getMaxPlayersForGame(dbRoom.GameType),
			rematchRequests: make(map[uuid.UUID]bool), // Initialize rematch map
		}
		m.rooms[room.ID] = room
		log.Printf("Activated room %s in memory.", room.ID)
	}
	m.mu.Unlock()

	log.Printf("🔄 handleJoinRoom: About to call room.addClient")
	room.addClient(client)
	log.Printf("✅ handleJoinRoom: Completed successfully")
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

	// Broadcast state update after every rematch request
	r.broadcastRoomState()

	if allPlayersRequestedRematch {
		r.mu.Lock()
		r.Game.Reset()

		// Swap roles for the new game
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

		// Broadcast the new game state after reset
		go r.broadcastRoomState()
	}
}
