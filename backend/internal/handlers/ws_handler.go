package handlers

import (
	"log"

	"github.com/DCCXXV/twoplayers/backend/internal/realtime"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	manager *realtime.Manager
}

func NewWebSocketHandler(m *realtime.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		manager: m,
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	err := h.manager.ServeWebSocket(c.Writer, c.Request)
	if err != nil {
		log.Printf("Error serving websocket: %v", err)
	}
}
