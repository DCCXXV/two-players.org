package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DCCXXV/twoplayers/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabaseConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	log.Println("Attempting to connect to database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return pool, nil
}
