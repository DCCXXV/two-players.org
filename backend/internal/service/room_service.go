package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRoomNotFound = errors.New("room not found")

type JoinRoomInput struct {
	RoomID      uuid.UUID
	DisplayName string
}

type JoinRoomResult struct {
	Success bool
	Role    string // "player_0", "player_1", "spectator"
	Message string
}

type RoomService interface {
	CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error)
	GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error)
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
	ListPublicRooms(ctx context.Context) ([]db.Room, error)
	ListPublicRoomsWithPlayers(ctx context.Context, gameType string, limit, offset int32) ([]db.ListPublicRoomsWithPlayersRow, error)
	JoinRoom(ctx context.Context, input JoinRoomInput) (*JoinRoomResult, error)
}

type CreateRoomParams struct {
	Name            string
	GameType        string
	HostDisplayName string
	IsPrivate       bool
	GameOptions     []byte
}

type roomService struct {
	queries *db.Queries
	db      *pgxpool.Pool
}

func NewRoomService(queries *db.Queries, db *pgxpool.Pool) RoomService {
	return &roomService{
		queries: queries,
		db:      db,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error) {
	return s.queries.CreateRoom(ctx, db.CreateRoomParams{
		Name:            params.Name,
		GameType:        params.GameType,
		HostDisplayName: params.HostDisplayName,
		IsPrivate:       params.IsPrivate,
		GameOptions:     params.GameOptions,
	})
}

func (s *roomService) GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error) {
	room, err := s.queries.GetRoomByID(ctx, pgtype.UUID{Bytes: roomID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Room{}, ErrRoomNotFound
		}
		return db.Room{}, err
	}
	return room, nil
}

func (s *roomService) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	return s.queries.DeleteRoom(ctx, pgtype.UUID{Bytes: roomID, Valid: true})
}

func (s *roomService) ListPublicRooms(ctx context.Context) ([]db.Room, error) {
	return s.queries.ListPublicRooms(ctx)
}

func (s *roomService) ListPublicRoomsWithPlayers(ctx context.Context, gameType string, limit, offset int32) ([]db.ListPublicRoomsWithPlayersRow, error) {
	return s.queries.ListPublicRoomsWithPlayers(ctx, db.ListPublicRoomsWithPlayersParams{
		GameType: gameType,
		Limit:    limit,
		Offset:   offset,
	})
}

func (s *roomService) JoinRoom(ctx context.Context, input JoinRoomInput) (*JoinRoomResult, error) {
	var result JoinRoomResult

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	players, err := qtx.GetPlayersByRoomID(ctx, pgtype.UUID{Bytes: input.RoomID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("could not get players in room: %w", err)
	}

	for _, p := range players {
		if p.PlayerDisplayName == input.DisplayName {
			return &JoinRoomResult{
				Success: true,
				Role:    fmt.Sprintf("player_%d", p.PlayerOrder),
				Message: "You are already in this game.",
			}, nil
		}
	}

	if len(players) < 2 {
		playerOrder := len(players)

		_, err := qtx.CreatePlayer(ctx, db.CreatePlayerParams{
			RoomID:            pgtype.UUID{Bytes: input.RoomID, Valid: true},
			PlayerDisplayName: input.DisplayName,
			PlayerOrder:       int16(playerOrder),
		})

		if err != nil {
			return nil, fmt.Errorf("could not add player to game: %w", err)
		}

		result.Success = true
		result.Role = fmt.Sprintf("player_%d", playerOrder)
		result.Message = "You have joined the game."
		log.Printf("Player %s joined room %s as player %d", input.DisplayName, input.RoomID, playerOrder)

	} else {
		result.Success = true
		result.Role = "spectator"
		result.Message = "The room is full. You are joining as a spectator."
		log.Printf("Player %s joined room %s as spectator", input.DisplayName, input.RoomID)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &result, nil
}
