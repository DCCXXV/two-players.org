-- name: CreatePlayer :one
-- Links an active connection (identified by player_display_name) to a specific room (room_id)
-- with a designated player_order (0 or 1).
-- Assumes the player_display_name already exists in the active_connections table.
INSERT INTO players (
    room_id,
    player_display_name,
    player_order
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetPlayersByRoomID :many
-- Retrieves all players associated with a specific room, ordered by their turn.
SELECT * FROM players
WHERE room_id = $1
ORDER BY player_order;
