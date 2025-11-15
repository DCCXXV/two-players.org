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
	UpdateConnectionName(ctx context.Context, oldName, newName string) error
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
	conn, err := s.queries.GetActiveConnection(ctx, params.DisplayName)
	if err == nil {
		return conn, nil
	}

	if err.Error() != "no rows in result set" {
		return db.ActiveConnection{}, fmt.Errorf("failed to get active connection: %w", err)
	}

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

func (s *connectionService) UpdateConnectionName(ctx context.Context, oldName, newName string) error {
	rowsAffected, err := s.queries.UpdateActiveConnectionName(ctx, db.UpdateActiveConnectionNameParams{
		DisplayName:   newName,
		DisplayName_2: oldName,
	})
	if err != nil {
		return fmt.Errorf("failed to update connection name: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("connection not found: %s", oldName)
	}
	return nil
}
