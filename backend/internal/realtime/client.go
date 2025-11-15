package realtime

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn        *websocket.Conn
	manager     *Manager
	id          uuid.UUID
	displayName string
	currentRoom *Room
	send        chan []byte
	role        string
	joinedAt    time.Time
}

func (c *Client) readPump() {
	defer func() {
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
				c.manager.logger.Error("WebSocket unexpected close", "client_id", c.id, "display_name", c.displayName, "error", err)
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
		case "chat_message":
			if c.currentRoom != nil {
				c.currentRoom.handleChatMessage(c, msg.Payload)
			} else {
				c.sendError("Not in a room.")
			}
		case "update_display_name":
			c.manager.handleUpdateDisplayName(c, msg.Payload)
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.manager.logger.Error("Failed to write message", "client_id", c.id, "error", err)
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
		c.manager.logger.Error("Failed to marshal message", "type", msgType, "client_id", c.id, "error", err)
		return
	}
	select {
	case c.send <- jsonData:
	default:
		c.manager.logger.Warn("Client send channel full, dropping message", "display_name", c.displayName, "message_type", msgType)
	}
}

func (c *Client) sendError(message string) {
	c.sendMessage("error", ErrorPayload{Message: message})
}
