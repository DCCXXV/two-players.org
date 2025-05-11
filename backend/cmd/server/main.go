package main

import (
	"log"
	"net/http"
	"os"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/DCCXXV/twoplayers/backend/internal/database"
	"github.com/DCCXXV/twoplayers/backend/internal/handlers"
	appLogger "github.com/DCCXXV/twoplayers/backend/internal/logger"
	"github.com/DCCXXV/twoplayers/backend/internal/realtime"
	"github.com/DCCXXV/twoplayers/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func ManualCorsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		if origin == "" {
			allowed = true
		} else {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")                                                                                                                                 // Ajusta según necesites
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With") // Cabeceras permitidas
		} else {
			log.Printf("WARN: Manual CORS Middleware: Denied origin: %s", origin)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if c.Request.Method == http.MethodOptions {
			log.Printf("DEBUG: Manual CORS Middleware: Handling OPTIONS request for origin: %s", origin)
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		if allowed {
			c.Request.Header.Del("Origin")
			log.Printf("DEBUG: Manual CORS Middleware: Deleted Origin header before passing to next handler.")
		}

		c.Next()
	}
}

func main() {
	appLogger.Init()
	log := appLogger.Get()

	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("FATAL: Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// 2. Establish Database Connection
	pool, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Error("FATAL: Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Initialize sqlc Queries
	queries := db.New(pool)

	// 4. Initialize Services
	roomService := service.NewRoomService(queries)
	connectionService := service.NewConnectionService(queries)
	playerService := service.NewPlayerService(queries)

	// 5. Initialize Realtime Manager
	rtManager, err := realtime.NewManager(cfg, connectionService, roomService, playerService)
	if err != nil {
		log.Error("FATAL: Failed to initialize realtime manager", "error", err)
		os.Exit(1)
	}

	// 6. Initialize Gin Router
	router := gin.Default()

	router.Use(ManualCorsMiddleware(cfg.AllowedOrigins))

	// 7. Setup Handlers
	httpHandler := handlers.NewHTTPHandler(roomService)
	wsHandler := handlers.NewWebSocketHandler(rtManager)

	// 8. Register Routes
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/rooms", httpHandler.CreateRoom)
		apiV1.GET("/rooms/:roomId", httpHandler.GetRoom)
		apiV1.GET("/rooms", httpHandler.ListPublicRooms)
		apiV1.DELETE("/rooms", httpHandler.DeleteRoom)
	}

	router.GET("/ws", wsHandler.HandleConnection)

	// 9. Start Server
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
