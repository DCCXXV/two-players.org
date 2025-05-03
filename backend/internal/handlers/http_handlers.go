package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HTTPHandler struct {
	roomService service.RoomService
}

func NewHTTPHandler(rs service.RoomService) *HTTPHandler {
	return &HTTPHandler{
		roomService: rs,
	}
}

// Request/Response Structs

type CreateRoomRequest struct {
	Name        string           `json:"name" binding:"required"`
	GameType    string           `json:"game_type" binding:"required"`
	IsPrivate   bool             `json:"is_private"`
	GameOptions *json.RawMessage `json:"game_options,omitempty"`
}

// Handlers

func (h *HTTPHandler) CreateRoom(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateRoomRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Printf("WARN: Bad request to create room: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	serviceParams := service.CreateRoomParams{
		Name:      req.Name,
		GameType:  req.GameType,
		IsPrivate: req.IsPrivate,
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

	log.Printf("INFO: Room created successfully: ID=%s, Name=%s, Type=%s", createdRoom.ID, createdRoom.Name, createdRoom.GameType)
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
		log.Printf("ERROR: Failed to get room %s: %v", roomIDStr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room"})
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
	rooms, err := h.roomService.ListPublicRooms(ctx)

	if err != nil {
		log.Printf("ERROR: Failed to list public rooms: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}

	if rooms == nil {
		rooms = []db.Room{}
	}

	c.JSON(http.StatusOK, rooms)
}
