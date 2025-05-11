-- name: CreateActiveConnection :one
-- Registers a new active connection with a unique display name.
-- Fails if the display_name is already taken (due to PRIMARY KEY constraint).
INSERT INTO active_connections (
    display_name
    -- last_seen defaults to NOW()
    -- status defaults to 'lobby'
    -- current_room_id defaults to NULL
) VALUES (
    $1
)
RETURNING *;

-- name: UpdateConnectionStatusAndRoom :one
-- Updates the status and current_room_id for an active connection.
-- e.g: when a player joins or leaves a room.
UPDATE active_connections
SET
    status = $2,
    current_room_id = $3,
    last_seen = NOW()
WHERE
    display_name = $1
RETURNING *;

-- name: UpdateActiveConnectionName :execrows
UPDATE active_connections
SET display_name = $1
WHERE display_name = $2
RETURNING display_name;

-- name: UpdateConnectionLastSeen :exec
-- Updates the last_seen timestamp for a connection (heartbeat).
UPDATE active_connections
SET last_seen = NOW()
WHERE display_name = $1;

-- name: DeleteActiveConnection :exec
-- Removes an active connection record (e.g., on disconnect).
-- ON DELETE CASCADE on players table will remove associated player records.
DELETE FROM active_connections
WHERE display_name = $1;

-- name: GetActiveConnection :one
-- Retrieves an active connection by display name.
SELECT * FROM active_connections
WHERE display_name = $1;

-- name: ListActiveLobbyUsers :many
-- Lists users currently in the lobby state.
SELECT display_name FROM active_connections
WHERE status = 'lobby'
ORDER BY display_name;

-- name: FindStaleConnections :many
-- Finds connections that haven't been seen recently (for cleanup).
SELECT display_name FROM active_connections
WHERE last_seen < $1; -- $1 would be a timestamp like NOW() - INTERVAL '5 minutes'
