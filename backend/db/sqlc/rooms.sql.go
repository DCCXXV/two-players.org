// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: rooms.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (
    name,
    game_type,
    host_display_name,
    game_options,
    is_private
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, name, game_type, game_options, is_private, created_at, host_display_name
`

type CreateRoomParams struct {
	Name            string `json:"name"`
	GameType        string `json:"game_type"`
	HostDisplayName string `json:"host_display_name"`
	GameOptions     []byte `json:"game_options"`
	IsPrivate       bool   `json:"is_private"`
}

// Create a new game room and return the created room record.
func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRow(ctx, createRoom,
		arg.Name,
		arg.GameType,
		arg.HostDisplayName,
		arg.GameOptions,
		arg.IsPrivate,
	)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.GameType,
		&i.GameOptions,
		&i.IsPrivate,
		&i.CreatedAt,
		&i.HostDisplayName,
	)
	return i, err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1
`

// Delete a room by its ID.
// Note: ON DELETE CASCADE on players table will handle removing associated players.
func (q *Queries) DeleteRoom(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteRoom, id)
	return err
}

const getRoomByID = `-- name: GetRoomByID :one
SELECT id, name, game_type, game_options, is_private, created_at, host_display_name FROM rooms
WHERE id = $1
LIMIT 1
`

// Retrieve a specific room by its unique ID.
func (q *Queries) GetRoomByID(ctx context.Context, id pgtype.UUID) (Room, error) {
	row := q.db.QueryRow(ctx, getRoomByID, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.GameType,
		&i.GameOptions,
		&i.IsPrivate,
		&i.CreatedAt,
		&i.HostDisplayName,
	)
	return i, err
}

const listPublicRooms = `-- name: ListPublicRooms :many
SELECT id, name, game_type, game_options, is_private, created_at, host_display_name FROM rooms
WHERE is_private = FALSE
ORDER BY created_at DESC
`

// Retrieve all rooms that are not private, ordered by creation time descending.
// Useful for a public lobby list.
func (q *Queries) ListPublicRooms(ctx context.Context) ([]Room, error) {
	rows, err := q.db.Query(ctx, listPublicRooms)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.GameType,
			&i.GameOptions,
			&i.IsPrivate,
			&i.CreatedAt,
			&i.HostDisplayName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPublicRoomsWithPlayers = `-- name: ListPublicRoomsWithPlayers :many
SELECT
    r.id,
    r.name,
    r.is_private,
    p1.player_display_name AS created_by,
    p2.player_display_name AS other_player
FROM rooms r
LEFT JOIN players p1 ON p1.room_id = r.id AND p1.player_order = 0
LEFT JOIN players p2 ON p2.room_id = r.id AND p2.player_order = 1
WHERE r.is_private = FALSE AND r.game_type = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3
`

type ListPublicRoomsWithPlayersParams struct {
	GameType string `json:"game_type"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

type ListPublicRoomsWithPlayersRow struct {
	ID          pgtype.UUID `json:"id"`
	Name        string      `json:"name"`
	IsPrivate   bool        `json:"is_private"`
	CreatedBy   pgtype.Text `json:"created_by"`
	OtherPlayer pgtype.Text `json:"other_player"`
}

// Retrieve all rooms that are not private with player info, ordered by creation time descending.
func (q *Queries) ListPublicRoomsWithPlayers(ctx context.Context, arg ListPublicRoomsWithPlayersParams) ([]ListPublicRoomsWithPlayersRow, error) {
	rows, err := q.db.Query(ctx, listPublicRoomsWithPlayers, arg.GameType, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPublicRoomsWithPlayersRow
	for rows.Next() {
		var i ListPublicRoomsWithPlayersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsPrivate,
			&i.CreatedBy,
			&i.OtherPlayer,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoomsByGameType = `-- name: ListRoomsByGameType :many
SELECT id, name, game_type, game_options, is_private, created_at, host_display_name FROM rooms
WHERE game_type = $1 AND is_private = FALSE
ORDER BY created_at DESC
`

// Retrieve all public rooms for a specific game type.
func (q *Queries) ListRoomsByGameType(ctx context.Context, gameType string) ([]Room, error) {
	rows, err := q.db.Query(ctx, listRoomsByGameType, gameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.GameType,
			&i.GameOptions,
			&i.IsPrivate,
			&i.CreatedAt,
			&i.HostDisplayName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
