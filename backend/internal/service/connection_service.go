package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
)

var ErrDisplayNameTaken = errors.New("display name already taken")
var namePrefixes = []string{"Alice", "Bob"}

const maxGenerateNameRetries = 5

func generateGuestName() string {
	prefix := namePrefixes[rand.Intn(len(namePrefixes))]
	randomNumber := rand.Intn(90000) + 10000
	return fmt.Sprintf("%s%d", prefix, randomNumber)
}

type ConnectionService interface {
	CreateConnection(ctx context.Context, params CreateConnectionParams) (db.ActiveConnection, error)
	DeleteConnection(ctx context.Context, displayName string) error
	ListActiveConnections(ctx context.Context) ([]db.ListActiveConnectionsRow, error)
}

type CreateConnectionParams struct {
	DisplayName string
}

type connectionService struct {
	queries db.Querier
}

func NewConnectionService(queries db.Querier) ConnectionService {
	return &connectionService{queries: queries}
}

func (s *connectionService) CreateConnection(ctx context.Context, params CreateConnectionParams) (db.ActiveConnection, error) {
	// Try to get the connection first
	conn, err := s.queries.GetActiveConnection(ctx, params.DisplayName)
	if err == nil {
		// Connection already exists, return it
		return conn, nil
	}

	// If the error is not sql.ErrNoRows, then it's a real error
	if err.Error() != "no rows in result set" {
		return db.ActiveConnection{}, fmt.Errorf("failed to get active connection: %w", err)
	}

	// Connection does not exist, create a new one
	newConn, err := s.queries.CreateActiveConnection(ctx, params.DisplayName)
	if err != nil {
		return db.ActiveConnection{}, fmt.Errorf("failed to create active connection: %w", err)
	}
	return newConn, nil
}

func (s *connectionService) DeleteConnection(ctx context.Context, displayName string) error {
	return s.queries.DeleteActiveConnection(ctx, displayName)
}

func (s *connectionService) ListActiveConnections(ctx context.Context) ([]db.ListActiveConnectionsRow, error) {
	return s.queries.ListActiveConnections(ctx)
}
