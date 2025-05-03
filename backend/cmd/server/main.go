package main

import (
	"log"
	"net/http" // Needed for splitting origins
	"time"     // Needed for CORS MaxAge

	// Adjust import paths based on your go.mod module name
	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/DCCXXV/twoplayers/backend/internal/database"
	"github.com/DCCXXV/twoplayers/backend/internal/handlers"
	"github.com/DCCXXV/twoplayers/backend/internal/realtime"
	"github.com/DCCXXV/twoplayers/backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// socketio "github.com/googollee/go-socket.io" // Keep if needed by realtime pkg
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v", err)
	}

	// 2. Establish Database Connection
	pool, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// 3. Initialize sqlc Queries
	queries := db.New(pool)

	// 4. Initialize Services
	roomService := service.NewRoomService(queries)
	connectionService := service.NewConnectionService(queries)
	playerService := service.NewPlayerService(queries)

	// 5. Initialize Realtime Manager
	rtManager, err := realtime.NewManager(connectionService, playerService, roomService)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize realtime manager: %v", err)
	}
	go rtManager.RunCleanupTask()

	// 6. Initialize Gin Router
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// 7. Setup Handlers
	httpHandler := handlers.NewHTTPHandler(roomService)

	// 8. Register Routes
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/rooms", httpHandler.CreateRoom)
		apiV1.GET("/rooms/:roomId", httpHandler.GetRoom)
		apiV1.DELETE("/rooms", httpHandler.ListPublicRooms)
		apiV1.GET("/rooms", httpHandler.ListPublicRooms)
	}

	socketIOServer := rtManager.GetServer()

	socketIOServer.SetAllowOriginFunc(func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, allowed := range cfg.AllowedOrigins {
			if origin == allowed {
				return true
			}
		}
		log.Printf("CORS check failed for Socket.IO origin: %s", origin)
		return false
	})

	socketIORoute := router.Group("/socket.io")
	{
		socketIORoute.GET("/*any", gin.WrapH(socketIOServer))
		socketIORoute.POST("/*any", gin.WrapH(socketIOServer))
	}

	// 9. Start Server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}
}
