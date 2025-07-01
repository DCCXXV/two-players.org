package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	roomService      service.RoomService
	playerService    service.PlayerService
	connectionService service.ConnectionService
}

func NewHTTPHandler(rs service.RoomService, ps service.PlayerService, cs service.ConnectionService) *HTTPHandler {
	return &HTTPHandler{
		roomService:      rs,
		playerService:    ps,
		connectionService: cs,
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
		log.Printf("WARN: Bad request to create room: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	displayName := c.GetHeader("X-Display-Name")
	if displayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing display name"})
		return
	}

	// Create active connection for the host
	_, err := h.connectionService.CreateConnection(ctx, service.CreateConnectionParams{DisplayName: displayName})
	if err != nil {
		if err == service.ErrDisplayNameTaken {
			c.JSON(http.StatusConflict, gin.H{"error": "Display name already taken"})
			return
		}
		log.Printf("ERROR: Failed to create active connection for host %s: %v", displayName, err)
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
		log.Printf("ERROR: Failed to create room via service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	playerParams := service.CreatePlayerParams{
		RoomID:            createdRoom.ID,
		PlayerDisplayName: displayName,
		PlayerOrder:       0,
	}
	_, err = h.playerService.CreatePlayer(ctx, playerParams)
	if err != nil {
		log.Printf("ERROR: Failed to create player 1 for room %s: %v", createdRoom.ID, err)
		delErr := h.roomService.DeleteRoom(ctx, createdRoom.ID.Bytes)
		if delErr != nil {
			log.Printf("ERROR: Failed to delete room %s after player creation failed: %v", createdRoom.ID, delErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create player for room"})
		return
	}

	log.Printf("INFO: Room and player 1 created successfully: RoomID=%s, Player=%s", createdRoom.ID, displayName)
	c.JSON(http.StatusCreated, createdRoom)
}

func (h *HTTPHandler) GetRoom(c *gin.Context) {
	ctx := c.Request.Context()
	roomIDStr := c.Param("roomId")

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		log.Printf("WARN: Invalid room ID format in URL: %s, error: %v", roomIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	room, err := h.roomService.GetRoomByID(ctx, roomID)
	if err != nil {
		if err == service.ErrRoomNotFound { // Check for specific error
			log.Printf("INFO: Room not found: %s", roomIDStr)
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		} else {
			log.Printf("ERROR: Failed to get room %s: %v", roomIDStr, err)
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
		log.Printf("WARN: Invalid room ID format in URL for delete: %s, error: %v", roomIDStr, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID format"})
		return
	}

	err = h.roomService.DeleteRoom(ctx, roomID)
	if err != nil {
		log.Printf("ERROR: Failed to delete room %s: %v", roomIDStr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	log.Printf("INFO: Room deleted successfully (or did not exist): ID=%s", roomIDStr)
	c.Status(http.StatusNoContent)
}

func (h *HTTPHandler) ListPublicRooms(c *gin.Context) {
	ctx := c.Request.Context()

	limit := int32(20)
	offset := int32(0)
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	rooms, err := h.roomService.ListPublicRoomsWithPlayers(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	type RoomPreview struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		IsPrivate   bool    `json:"is_private"`
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
			IsPrivate:   r.IsPrivate,
			CreatedBy:   createdBy,
			OtherPlayer: otherPlayer,
		})
	}

	c.JSON(http.StatusOK, previews)
}
