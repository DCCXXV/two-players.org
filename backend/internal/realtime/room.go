package realtime

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

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
	r.mu.Lock()

	if _, ok := r.Clients[client.id]; ok {
		r.mu.Unlock()
		return
	}

	playerCount := r.getPlayerCountInternal()

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
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				r.manager.logger.Debug("Player already exists in room", "display_name", client.displayName, "room_id", r.ID)
			} else {
				r.manager.logger.Error("Failed to create player in DB", "display_name", client.displayName, "room_id", r.ID, "error", err)
				client.sendError("Failed to join room: could not save player data.")
				r.mu.Unlock()
				return
			}
		}
	}

	client.joinedAt = time.Now()
	client.currentRoom = r
	r.Clients[client.id] = client

	client.sendMessage("join_success", map[string]any{
		"roomId": r.ID.String(),
		"role":   client.role,
	})

	r.mu.Unlock()

	r.broadcastRoomState()
}

func (r *Room) removeClient(client *Client) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Clients[client.id]; !ok {
		return false
	}

	isHost := client.displayName == r.HostName

	delete(r.Clients, client.id)
	client.currentRoom = nil

	if isHost {
		for _, otherClient := range r.Clients {
			otherClient.sendMessage("room_closed", map[string]string{"message": "The host has left the room."})
			otherClient.currentRoom = nil
		}

		r.Clients = make(map[uuid.UUID]*Client)
		return true

	} else if len(r.Clients) == 0 {
		return true
	} else {
		go r.broadcastRoomState()
		return false
	}
}

func (r *Room) broadcastRoomState() {
	r.mu.RLock()
	players := r.getPlayersInternal()
	spectators := r.getSpectatorsInternal()

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

	r.broadcastMessage("game_state_update", roomState)
}

func (r *Room) broadcastMessage(msgType string, payload any) {
	message, err := createWebSocketMessage(msgType, payload)
	if err != nil {
		r.manager.logger.Error("Failed to create message", "type", msgType, "room_id", r.ID, "error", err)
		return
	}

	r.mu.RLock()
	clientsCopy := make([]*Client, 0, len(r.Clients))
	for _, client := range r.Clients {
		clientsCopy = append(clientsCopy, client)
	}
	r.mu.RUnlock()

	for _, client := range clientsCopy {
		select {
		case client.send <- message:
		default:
			r.manager.logger.Warn("Client send channel full in broadcast", "display_name", client.displayName, "room_id", r.ID, "message_type", msgType)
		}
	}
}
