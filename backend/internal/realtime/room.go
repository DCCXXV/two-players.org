package realtime

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// Room represents an active game room in memory.
type Room struct {
	ID              uuid.UUID
	GameType        string
	HostName        string
	Game            GameInstance
	manager         *Manager
	mu              sync.RWMutex
	Clients         map[uuid.UUID]*Client
	MaxPlayers      int
	rematchRequests map[uuid.UUID]bool
}

func (r *Room) getPlayersInternal() []*Client {
	var player0, player1 *Client
	for _, client := range r.Clients {
		if client.role == "player_0" {
			player0 = client
		} else if client.role == "player_1" {
			player1 = client
		}
	}

	var players []*Client
	if player0 != nil {
		players = append(players, player0)
	}
	if player1 != nil {
		players = append(players, player1)
	}
	return players
}

func (r *Room) getSpectatorsInternal() []*Client {
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

func (r *Room) addClient(client *Client) {
	log.Printf("🔄 addClient: Starting for client %s in room %s", client.displayName, r.ID)

	r.mu.Lock()

	if _, ok := r.Clients[client.id]; ok {
		log.Printf("⚠️  addClient: Client already in room")
		r.mu.Unlock()
		return
	}

	playerCount := r.getPlayerCountInternal()
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

	// If the client is a player, persist them in the DB before adding to the in-memory room.
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
			// If it's a unique constraint violation, it's not a fatal error here.
			// It means the player was already registered (e.g., the host rejoining).
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

	r.mu.Unlock()

	r.broadcastRoomState()
}

// removeClient removes a client from the room. It returns `true` if the room becomes empty and should be deleted.
// This function is designed to be called from `manager.unregisterClient` to avoid deadlocks.
func (r *Room) removeClient(client *Client) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Printf("Room %s: removeClient called for client %s (%s). Host: %s. Current clients in room: %d", r.ID, client.displayName, client.id, r.HostName, len(r.Clients))

	if _, ok := r.Clients[client.id]; !ok {
		log.Printf("Room %s: Client %s not found in room. Aborting removeClient.", r.ID, client.displayName)
		return false // The room should not be deleted
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

		r.Clients = make(map[uuid.UUID]*Client)
		log.Printf("Room %s: All clients notified and room client map cleared. Signalling for deletion.", r.ID)
		return true // Indicate that the room should be deleted

	} else if len(r.Clients) == 0 {
		log.Printf("Room %s: Last non-host client %s left. Room is now empty. Signalling for deletion.", r.ID, client.displayName)
		return true // Indicate that the room should be deleted
	} else {
		// A non-host client left, just update the state for the others.
		log.Printf("Room %s: Non-host client %s left. Broadcasting updated room state.", r.ID, client.displayName)
		go r.broadcastRoomState() // Use goroutine to avoid deadlock
		return false              // The room should NOT be deleted
	}
}

func (r *Room) broadcastRoomState() {
	log.Printf("🔄 broadcastRoomState: Starting for room %s", r.ID)

	r.mu.RLock()
	players := r.getPlayersInternal()
	spectators := r.getSpectatorsInternal()

	// Convert to names for the frontend
	playerNames := make([]string, len(players))
	for i, p := range players {
		playerNames[i] = p.displayName
	}

	spectatorNames := make([]string, len(spectators))
	for i, s := range spectators {
		spectatorNames[i] = s.displayName
	}

	gameState := r.Game.GetGameState()
	rematchCount := len(r.rematchRequests)
	r.mu.RUnlock()

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
		"rematchCount":   rematchCount,
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
