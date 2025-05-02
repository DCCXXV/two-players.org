-- name: CreateRoom :one
INSERT INTO rooms (id, game_type, is_private)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRooms :many
SELECT id, game_type, is_private, created_at
FROM rooms
ORDER BY created_at DESC;