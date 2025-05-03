-- name: CreateRoom :one
-- Create a new game room and return the created room record.
INSERT INTO rooms (
    name,
    game_type,
    game_options,
    is_private
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetRoomByID :one
-- Retrieve a specific room by its unique ID.
SELECT * FROM rooms
WHERE id = $1
LIMIT 1;

-- name: ListPublicRooms :many
-- Retrieve all rooms that are not private, ordered by creation time descending.
-- Useful for a public lobby list.
SELECT * FROM rooms
WHERE is_private = FALSE
ORDER BY created_at DESC;

-- name: ListRoomsByGameType :many
-- Retrieve all public rooms for a specific game type.
SELECT * FROM rooms
WHERE game_type = $1 AND is_private = FALSE
ORDER BY created_at DESC;

-- name: DeleteRoom :exec
-- Delete a room by its ID.
-- Note: ON DELETE CASCADE on players table will handle removing associated players.
DELETE FROM rooms
WHERE id = $1;