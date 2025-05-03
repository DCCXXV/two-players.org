package service

import (
	"context"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
)

type ConnectionService interface {
	CreateConnection(ctx context.Context, params CreateConnectionParams) (db.ActiveConnection, error)
}

type CreateConnectionParams struct {
	DisplayName string
}

type connectionService struct {
	queries db.Querier
}

func NewConnectionService(queries db.Querier) connectionService {
	return &connectionService{
		queries: queries,
	}
}
