package realtime

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a single user connected via WebSocket.
type Client struct {
	conn        *websocket.Conn
	manager     *Manager
	id          uuid.UUID
	displayName string
	currentRoom *Room
	send        chan []byte

	// State in the current room
	role     string // "player_0", "player_1", "spectator"
	joinedAt time.Time
}

func (c *Client) readPump() {
	defer func() {
		log.Printf("Realtime: Closing readPump for client %s (%s)", c.id, c.displayName)
		c.manager.unregisterClient(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

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
			// Handle game moves
			if c.currentRoom != nil {
				c.handleGameMove(msg.Payload)
			} else {
				c.sendError("Not in a room.")
			}
		case "rematch_request":
			if c.currentRoom != nil {
				c.currentRoom.handleRematch(c)
			} else {
				c.sendError("Not in a room.")
			}
		default:
			c.sendError(fmt.Sprintf("Unknown message type '%s'.", msg.Type))
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The manager closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Send each message in its own frame to prevent JSON concatenation.
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("ERROR: Failed to write message for client %s: %v", c.id, err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
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
	select {
	case c.send <- jsonData:
	default:
		log.Printf("⚠️  sendMessage: Client %s's send channel is full. Dropping message of type '%s'.", c.displayName, msgType)
	}
}

func (c *Client) sendError(message string) {
	c.sendMessage("error", ErrorPayload{Message: message})
}
