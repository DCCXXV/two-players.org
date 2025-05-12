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
	return s.queries.CreateActiveConnection(ctx, params.DisplayName)
}

func (s *connectionService) DeleteConnection(ctx context.Context, displayName string) error {
	return s.queries.DeleteActiveConnection(ctx, displayName)
}
