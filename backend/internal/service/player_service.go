package service

import (
	"context"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, params CreatePlayerParams) (db.Player, error)
	DeletePlayerByRoomAndName(ctx context.Context, roomID pgtype.UUID, playerDisplayName string) error
	DeletePlayersByRoomID(ctx context.Context, roomID pgtype.UUID) error
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

func (s *playerService) DeletePlayerByRoomAndName(ctx context.Context, roomID pgtype.UUID, playerDisplayName string) error {
	return s.queries.DeletePlayerByRoomAndName(ctx, db.DeletePlayerByRoomAndNameParams{
		RoomID:            roomID,
		PlayerDisplayName: playerDisplayName,
	})
}

func (s *playerService) DeletePlayersByRoomID(ctx context.Context, roomID pgtype.UUID) error {
	return s.queries.DeletePlayersByRoomID(ctx, roomID)
}
