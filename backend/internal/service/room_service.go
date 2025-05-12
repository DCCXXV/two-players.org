package service

import (
	"context"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomService interface {
	CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error)
	GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error)
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
	ListPublicRooms(ctx context.Context) ([]db.Room, error)
	ListPublicRoomsWithPlayers(ctx context.Context, limit, offset int32) ([]db.ListPublicRoomsWithPlayersRow, error)
}

type CreateRoomParams struct {
	Name        string
	GameType    string
	IsPrivate   bool
	GameOptions []byte
}

type roomService struct {
	queries db.Querier
}

func NewRoomService(queries db.Querier) RoomService {
	return &roomService{queries: queries}
}

func (s *roomService) CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error) {
	return s.queries.CreateRoom(ctx, db.CreateRoomParams{
		Name:        params.Name,
		GameType:    params.GameType,
		IsPrivate:   params.IsPrivate,
		GameOptions: params.GameOptions,
	})
}

func (s *roomService) GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error) {
	return s.queries.GetRoomByID(ctx, pgtype.UUID{Bytes: roomID, Valid: true})
}

func (s *roomService) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	return s.queries.DeleteRoom(ctx, pgtype.UUID{Bytes: roomID, Valid: true})
}

func (s *roomService) ListPublicRooms(ctx context.Context) ([]db.Room, error) {
	return s.queries.ListPublicRooms(ctx)
}

func (s *roomService) ListPublicRoomsWithPlayers(ctx context.Context, limit, offset int32) ([]db.ListPublicRoomsWithPlayersRow, error) {
	return s.queries.ListPublicRoomsWithPlayers(ctx, db.ListPublicRoomsWithPlayersParams{
		Limit:  limit,
		Offset: offset,
	})
}
