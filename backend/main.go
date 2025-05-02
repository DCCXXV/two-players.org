package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type CreateRoomRequest struct {
	GameType  string `json:"game_type" binding:"required"`
	IsPrivate bool   `json:"is_private"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file: ", err)
	}

	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		log.Fatal("failed to parse pgxpool config: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("db ping failed: ", err)
	}
	log.Println("Database connected succesfully")

	queries := db.New(pool)

	router := gin.Default()

	router.GET("/api/rooms", listRoomsHandler(queries))
	router.POST("/api/rooms", createRoomHandler(queries))

	port := ":8080"
	log.Printf("server starting on port %s", port)
	err = router.Run(port)
	if err != nil {
		log.Fatal("server failed to start:", err)
	}
}

func listRoomsHandler(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms, err := queries.ListRooms(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rooms"})
			return
		}
		c.JSON(http.StatusOK, rooms)
	}
}

func createRoomHandler(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody CreateRoomRequest

		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Printf("Bad request body: %v", err)
			return
		}

		if reqBody.GameType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "game_type is required"})
			return
		}

		roomID, err := uuid.NewRandom()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate room ID"})
			log.Printf("Error generating UUID: %v", err)
			return
		}

		var pgUUID pgtype.UUID
		copy(pgUUID.Bytes[:], roomID[:])
		pgUUID.Valid = true

		createdRoom, err := queries.CreateRoom(c.Request.Context(), db.CreateRoomParams{
			ID:        pgUUID,
			GameType:  reqBody.GameType,
			IsPrivate: reqBody.IsPrivate,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
			log.Printf("Error creating room in DB: %v", err)
			return
		}

		c.JSON(http.StatusCreated, createdRoom)
		log.Printf("Room created successfully: %v", createdRoom.ID)
	}
}
