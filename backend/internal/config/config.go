// Package config handles loading application configuration from environment
// variables and .env files.
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	ServerPort     string
	AllowedOrigins string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: Could not load .env file:", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
		log.Printf("Info: SERVER_PORT not set, defaulting to %s", serverPort)
	}

	if serverPort[0] != ':' {
		serverPort = ":" + serverPort
	}

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:5173"
		log.Printf("Info: ALLOWED_ORIGINS not set, defaulting to %s", allowedOriginsEnv)
	}

	cfg := &Config{
		DatabaseURL:    dbURL,
		ServerPort:     serverPort,
		AllowedOrigins: allowedOriginsEnv,
	}

	log.Printf("Configuration loaded: Port=%s, AllowedOrigins=%v", cfg.ServerPort, cfg.AllowedOrigins)
	return cfg, nil
}
