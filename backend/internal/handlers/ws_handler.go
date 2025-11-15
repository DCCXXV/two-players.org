package handlers

import (
	"log/slog"

	appLogger "github.com/DCCXXV/twoplayers/backend/internal/logger"
	"github.com/DCCXXV/twoplayers/backend/internal/realtime"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	manager *realtime.Manager
	logger  *slog.Logger
}

func NewWebSocketHandler(m *realtime.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		manager: m,
		logger:  appLogger.Get(),
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	err := h.manager.ServeWebSocket(c.Writer, c.Request)
	if err != nil {
		h.logger.Error("Error serving websocket", "error", err, "remote_addr", c.Request.RemoteAddr)
	}
}
