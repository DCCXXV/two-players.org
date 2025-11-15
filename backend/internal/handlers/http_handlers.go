package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	appLogger "github.com/DCCXXV/twoplayers/backend/internal/logger"
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	roomService       service.RoomService
	playerService     service.PlayerService
	connectionService service.ConnectionService
	logger            *slog.Logger
}

func NewHTTPHandler(rs service.RoomService, ps service.PlayerService, cs service.ConnectionService) *HTTPHandler {
	return &HTTPHandler{
		roomService:       rs,
		playerService:     ps,
		connectionService: cs,
		logger:            appLogger.Get(),
	}
}

type CreateRoomRequest struct {
	Name        string           `json:"name" binding:"required"`
	GameType    string           `json:"game_type" binding:"required"`
	IsPrivate   bool             `json:"is_private"`
	GameOptions *json.RawMessage `json:"game_options,omitempty"`
}

func (h *HTTPHandler) CreateRoom(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateRoomRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.logger.Warn("Bad request to create room", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	displayName := c.GetHeader("X-Display-Name")
	if displayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing display name"})
		return
	}

	_, err := h.connectionService.CreateConnection(ctx, service.CreateConnectionParams{DisplayName: displayName})
	if err != nil {
		if err == service.ErrDisplayNameTaken {
			c.JSON(http.StatusConflict, gin.H{"error": "Display name already taken"})
			return
		}
		h.logger.Error("Failed to create active connection for host", "display_name", displayName, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish host connection"})
		return
	}

	serviceParams := service.CreateRoomParams{
		Name:            req.Name,
		GameType:        req.GameType,
		IsPrivate:       req.IsPrivate,
		HostDisplayName: displayName,
	}
	if req.GameOptions != nil {
		serviceParams.GameOptions = *req.GameOptions
	}

	createdRoom, err := h.roomService.CreateRoom(ctx, serviceParams)
	if err != nil {
		h.logger.Error("Failed to create room via service", "error", err)
		if delConnErr := h.connectionService.DeleteConnection(ctx, displayName); delConnErr != nil {
			h.logger.Error("Failed to rollback connection", "display_name", displayName, "error", delConnErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	h.logger.Info("Room created successfully", "room_id", createdRoom.ID, "host", displayName)
	c.JSON(http.StatusCreated, createdRoom)
}

func (h *HTTPHandler) GetRoom(c *gin.Context) {
	ctx := c.Request.Context()
	roomIDStr := c.Param("roomId")

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		h.logger.Warn("Invalid room ID format in URL", "room_id", roomIDStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	room, err := h.roomService.GetRoomByID(ctx, roomID)
	if err != nil {
		if err == service.ErrRoomNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		} else {
			h.logger.Error("Failed to get room", "room_id", roomIDStr, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room"})
		}
		return
	}

	c.JSON(http.StatusOK, room)
}

func (h *HTTPHandler) DeleteRoom(c *gin.Context) {
	ctx := c.Request.Context()
	roomIDStr := c.Param("roomId")

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		h.logger.Warn("Invalid room ID format in URL for delete", "room_id", roomIDStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	err = h.roomService.DeleteRoom(ctx, roomID)
	if err != nil {
		h.logger.Error("Failed to delete room", "room_id", roomIDStr, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	h.logger.Info("Room deleted successfully", "room_id", roomIDStr)
	c.Status(http.StatusNoContent)
}

func (h *HTTPHandler) ListPublicRooms(c *gin.Context) {
	ctx := c.Request.Context()

	gameType := c.Query("game_type")
	if gameType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game_type query parameter is required"})
		return
	}

	limit := int32(20)
	offset := int32(0)
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	rooms, err := h.roomService.ListPublicRoomsWithPlayers(ctx, gameType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	type RoomPreview struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		GameType    string  `json:"game_type"`
		IsPrivate   bool    `json:"is_private"`
		CreatedAt   string  `json:"created_at"`
		CreatedBy   *string `json:"created_by"`
		OtherPlayer *string `json:"other_player"`
	}

	previews := make([]RoomPreview, 0, len(rooms))
	for _, r := range rooms {
		var createdBy *string
		if r.CreatedBy.Valid {
			createdBy = &r.CreatedBy.String
		}
		var otherPlayer *string
		if r.OtherPlayer.Valid {
			otherPlayer = &r.OtherPlayer.String
		}
		previews = append(previews, RoomPreview{
			ID:          r.ID.String(),
			Name:        r.Name,
			GameType:    r.GameType,
			IsPrivate:   r.IsPrivate,
			CreatedAt:   r.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			CreatedBy:   createdBy,
			OtherPlayer: otherPlayer,
		})
	}

	c.JSON(http.StatusOK, previews)
}

func (h *HTTPHandler) ListActiveConnections(c *gin.Context) {
	ctx := c.Request.Context()

	connections, err := h.connectionService.ListActiveConnections(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active connections"})
		return
	}

	c.JSON(http.StatusOK, connections)
}
