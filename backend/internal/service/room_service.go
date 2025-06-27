package service

import (
	"context"
	"fmt"
	"log"

	db "github.com/DCCXXV/twoplayers/backend/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool" // NUEVO: Necesario para las transacciones
)

// NUEVO: Struct para los parámetros de entrada de JoinRoom
type JoinRoomInput struct {
	RoomID      uuid.UUID
	DisplayName string
}

// NUEVO: Struct para el resultado de JoinRoom
type JoinRoomResult struct {
	Success bool
	Role    string // "player_0", "player_1", "spectator"
	Message string
}

// MODIFICADO: La interfaz ahora incluye JoinRoom
type RoomService interface {
	CreateRoom(ctx context.Context, params CreateRoomParams) (db.Room, error)
	GetRoomByID(ctx context.Context, roomID uuid.UUID) (db.Room, error)
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
	ListPublicRooms(ctx context.Context) ([]db.Room, error)
	ListPublicRoomsWithPlayers(ctx context.Context, limit, offset int32) ([]db.ListPublicRoomsWithPlayersRow, error)
	JoinRoom(ctx context.Context, input JoinRoomInput) (*JoinRoomResult, error) // NUEVO
}

type CreateRoomParams struct {
	Name            string
	GameType        string
	HostDisplayName string
	IsPrivate       bool
	GameOptions     []byte
}

// MODIFICADO: El struct ahora también guarda el pool de la BD para transacciones
type roomService struct {
	queries *db.Queries   // MODIFICADO: Usar el struct concreto para acceder a WithTx
	db      *pgxpool.Pool // NUEVO
}

// MODIFICADO: El constructor ahora acepta el pool de la BD
func NewRoomService(queries *db.Queries, db *pgxpool.Pool) RoomService { // MODIFICADO: Usar el struct concreto
	return &roomService{
		queries: queries,
		db:      db, // NUEVO
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

func (s *roomService) JoinRoom(ctx context.Context, input JoinRoomInput) (*JoinRoomResult, error) {
	var result JoinRoomResult

	// Iniciar transacción para garantizar la consistencia de los datos
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	// Asegurarse de que la transacción se anule si algo sale mal
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// 1. Obtener los jugadores actuales en la sala (dentro de la tx para bloquear)
	// Con pgx/v5 y sqlc, una consulta "many" que no devuelve filas resulta en un slice vacío y un error nil,
	// no en sql.ErrNoRows. Por lo tanto, solo necesitamos comprobar si hay un error real.
	players, err := qtx.GetPlayersByRoomID(ctx, pgtype.UUID{Bytes: input.RoomID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("could not get players in room: %w", err)
	}

	// 2. Verificar si el jugador ya está en la sala
	for _, p := range players {
		if p.PlayerDisplayName == input.DisplayName {
			return &JoinRoomResult{
				Success: true,
				Role:    fmt.Sprintf("player_%d", p.PlayerOrder),
				Message: "Ya estás en esta partida.",
			}, nil
		}
	}

	// 3. Decidir si se puede unir como jugador o debe ser espectador
	if len(players) < 2 {
		// Hay espacio, ¡unirse como jugador!
		playerOrder := len(players) // 0 para el primero, 1 para el segundo

		_, err := qtx.CreatePlayer(ctx, db.CreatePlayerParams{
			RoomID:            pgtype.UUID{Bytes: input.RoomID, Valid: true},
			PlayerDisplayName: input.DisplayName,
			PlayerOrder:       int16(playerOrder),
		})

		if err != nil {
			// Si falla, podría ser por una race condition que el constraint UNIQUE ha detenido.
			return nil, fmt.Errorf("could not add player to game: %w", err)
		}

		result.Success = true
		result.Role = fmt.Sprintf("player_%d", playerOrder)
		result.Message = "Te has unido a la partida."
		log.Printf("Player %s joined room %s as player %d", input.DisplayName, input.RoomID, playerOrder)

	} else {
		// La sala está llena, unirse como espectador.
		result.Success = true
		result.Role = "spectator"
		result.Message = "La sala está llena. Entras como espectador."
		log.Printf("Player %s joined room %s as spectator", input.DisplayName, input.RoomID)
	}

	// 4. Si todo fue bien, hacer commit de la transacción
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &result, nil
}
