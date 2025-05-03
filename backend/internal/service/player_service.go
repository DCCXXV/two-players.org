package service

import (
	"context"
	"fmt"
	"log"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, params CreatePlayerParams) (db.Player, error)
}

type CreatePlayerParams struct {
	roomID            pgtype.UUID
	playerDisplayName string
	playerOrder       int16
}

type playerService struct {
	queries db.Querier
}

func NewPlayerService(queries db.Querier) PlayerService {
	return &playerService{
		queries: queries,
	}
}

func (s *playerService) CreatePlayer(ctx context.Context, params CreatePlayerParams) (db.Player, error) {
	log.Print("Service: Attempting to create player")

	dbParams := db.CreatePlayerParams{
		RoomID:            params.roomID,
		PlayerDisplayName: params.playerDisplayName,
		PlayerOrder:       params.playerOrder,
	}

	createdPlayer, err := s.queries.CreatePlayer(ctx, dbParams)
	if err != nil {
		log.Printf("ERROR: Service failed to create player in DB: %v", err)
		return db.Player{}, fmt.Errorf("database error creating player: %w", err)
	}

	log.Printf("Service: Player created successfully with ID: %s", createdPlayer.ID)
	return createdPlayer, nil
}
