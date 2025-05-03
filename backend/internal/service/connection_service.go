package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
)

var ErrDisplayNameTaken = errors.New("display name already taken")
var namePrefixes = []string{"Alice", "Bob"}

const maxGenerateNameRetries = 5

func generateGuestName() string {
	prefixIndex := rand.Intn(len(namePrefixes))
	prefix := namePrefixes[prefixIndex]
	randomNumber := rand.Intn(90000) + 10000
	return fmt.Sprintf("%s%d", prefix, randomNumber)
}

type ConnectionService interface {
	CreateConnection(ctx context.Context, params CreateConnectionParams) (db.ActiveConnection, error)
	DeleteConnection(ctx context.Context, displayName string) error
}

type CreateConnectionParams struct {
	DisplayName string
}

type connectionService struct {
	queries db.Querier
}

func NewConnectionService(queries db.Querier) ConnectionService {
	return &connectionService{
		queries: queries,
	}
}

func (s *connectionService) CreateConnection(ctx context.Context, params CreateConnectionParams) (db.ActiveConnection, error) {
	log.Printf("Service: Attempting to create a connection with DisplayName=%s", params.DisplayName)

	activeConnection, err := s.queries.CreateActiveConnection(ctx, params.DisplayName)
	if err != nil {
		log.Printf("ERROR: Service failed to create connection in DB: %v", err)
		return db.ActiveConnection{}, fmt.Errorf("database error creating connection: %w", err)
	}

	log.Printf("Service: Connection created succesfully")
	return activeConnection, nil
}

func (s *connectionService) DeleteConnection(ctx context.Context, displayName string) error {
	log.Printf("Service: Attempting to create a connection with DisplayName=%s", displayName)

	err := s.queries.DeleteActiveConnection(ctx, displayName)
	if err != nil {
		log.Printf("ERROR: Service failed to delete connection %s: %v", displayName, err)
		return fmt.Errorf("database error deleting room: %w", err)
	}

	log.Printf("Service: Connection deleted (or did not exist) with ID: %s", displayName)
	return nil
}
