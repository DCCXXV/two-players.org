package main

import (
	"net/http"
	"os"
	"strings"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/DCCXXV/twoplayers/backend/internal/database"
	_ "github.com/DCCXXV/twoplayers/backend/internal/games"
	"github.com/DCCXXV/twoplayers/backend/internal/handlers"
	appLogger "github.com/DCCXXV/twoplayers/backend/internal/logger"
	"github.com/DCCXXV/twoplayers/backend/internal/realtime"
	"github.com/DCCXXV/twoplayers/backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	appLogger.Init()
	log := appLogger.Get()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("FATAL: Failed to load configuration", "error", err)
		os.Exit(1)
	}

	pool, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Error("FATAL: Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)

	roomService := service.NewRoomService(queries, pool)
	connectionService := service.NewConnectionService(queries)
	playerService := service.NewPlayerService(queries)

	rtManager, err := realtime.NewManager(cfg, connectionService, roomService, playerService)
	if err != nil {
		log.Error("FATAL: Failed to initialize realtime manager", "error", err)
		os.Exit(1)
	}

	router := gin.Default()
	allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Display-Name"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))
	router.Use(cors.New(corsConfig))

	httpHandler := handlers.NewHTTPHandler(roomService, playerService, connectionService)
	wsHandler := handlers.NewWebSocketHandler(rtManager)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/rooms", httpHandler.CreateRoom)
		apiV1.GET("/rooms/:roomId", httpHandler.GetRoom)
		apiV1.GET("/rooms", httpHandler.ListPublicRooms)
		apiV1.GET("/connections", httpHandler.ListActiveConnections)
		apiV1.DELETE("/rooms", httpHandler.DeleteRoom)
	}

	router.GET("/ws", wsHandler.HandleConnection)

	log.Info("Server starting on port", "port", cfg.ServerPort)
	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error("FATAL: Server failed to start", "error", err)
		os.Exit(1)
	}
}
