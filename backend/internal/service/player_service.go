package service

import (
	"context"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, params CreatePlayerParams) (db.Player, error)
}

type CreatePlayerParams struct {
	RoomID            pgtype.UUID
	PlayerDisplayName string
	PlayerOrder       int16
}

type playerService struct {
	queries db.Querier
}

func NewPlayerService(queries db.Querier) PlayerService {
	return &playerService{queries: queries}
}

func (s *playerService) CreatePlayer(ctx context.Context, params CreatePlayerParams) (db.Player, error) {
	return s.queries.CreatePlayer(ctx, db.CreatePlayerParams{
		RoomID:            params.RoomID,
		PlayerDisplayName: params.PlayerDisplayName,
		PlayerOrder:       params.PlayerOrder,
	})
}
