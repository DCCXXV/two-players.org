package service

import (
	"context"
	"fmt"
	"log"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoomService interface {
	CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error)
	GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error)
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
	ListPublicRooms(ctx context.Context) ([]db.Room, error)
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
	return &roomService{
		queries: queries,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error) {
	log.Printf("Service: Attempting to create room with GameType=%s", params.GameType)

	dbParams := db.CreateRoomParams{
		Name:        params.Name,
		GameType:    params.GameType,
		IsPrivate:   params.IsPrivate,
		GameOptions: params.GameOptions,
	}

	createdRoom, err := s.queries.CreateRoom(ctx, dbParams)
	if err != nil {
		log.Printf("ERROR: Service failed to create room in DB: %v", err)
		return db.Room{}, fmt.Errorf("database error creating room: %w", err)
	}

	log.Printf("Service: Room created successfully with ID: %s", createdRoom.ID)
	return createdRoom, nil
}

func (s *roomService) GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error) {
	log.Printf("Service: Getting room by ID: %s", roomID.String())

	pgUUID := pgtype.UUID{Bytes: roomID, Valid: true}
	room, err := s.queries.GetRoomByID(ctx, pgUUID)
	if err != nil {
		log.Printf("ERROR: Service failed to get room by ID %s: %v", roomID.String(), err)
		return db.Room{}, fmt.Errorf("database error getting room by ID: %w", err)
	}

	log.Printf("Service: Found room: %s", roomID.String())
	return room, nil
}

func (s *roomService) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	log.Printf("Service: Attempting to delete room by ID: %s", roomID.String())

	pgUUID := pgtype.UUID{Bytes: roomID, Valid: true}
	err := s.queries.DeleteRoom(ctx, pgUUID)
	if err != nil {
		log.Printf("ERROR: Service failed to delete room %s: %v", roomID.String(), err)
		return fmt.Errorf("database error deleting room: %w", err)
	}

	log.Printf("Service: Room deleted (or did not exist) with ID: %s", roomID.String())

	return nil
}
func (s *roomService) ListPublicRooms(ctx context.Context) ([]db.Room, error) {
	log.Println("Service: Listing public rooms")

	rooms, err := s.queries.ListPublicRooms(ctx)
	if err != nil {
		log.Printf("ERROR: Service failed to list public rooms: %v", err)
		return nil, fmt.Errorf("database error listing public rooms: %w", err)
	}

	log.Printf("Service: Found %d public rooms", len(rooms))
	return rooms, nil
}
