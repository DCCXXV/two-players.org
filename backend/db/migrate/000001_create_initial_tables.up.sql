--  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- -----------------------------------------------------
-- Table `rooms`
-- Stores information about each game room.
-- -----------------------------------------------------
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    game_type VARCHAR(50) NOT NULL,
    game_options JSONB,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rooms_game_type ON rooms(game_type);

-- -----------------------------------------------------
-- Table `active_connections`
-- Tracks currently connected users/sessions and their chosen display names.
-- Enforces global uniqueness of display names for active sessions.
-- -----------------------------------------------------
CREATE TABLE active_connections (
    display_name VARCHAR(50) PRIMARY KEY,
    last_seen TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'lobby' CHECK (status IN ('lobby', 'in_game', 'looking')),
    current_room_id UUID NULL,

    CONSTRAINT fk_current_room
        FOREIGN KEY(current_room_id)
        REFERENCES rooms(id)
        ON DELETE SET NULL -- Set the room ID to NULL if the room is deleted
);

CREATE INDEX idx_active_connections_last_seen ON active_connections(last_seen);
CREATE INDEX idx_active_connections_status ON active_connections(status);
CREATE INDEX idx_active_connections_current_room_id ON active_connections(current_room_id) WHERE current_room_id IS NOT NULL;

-- -----------------------------------------------------
-- Table `players`
-- Links an active connection (user) to a specific room, assigning their role (player order).
-- Represents a player actively participating in a game instance.
-- -----------------------------------------------------
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL,
    player_display_name VARCHAR(50) NOT NULL,
    player_order SMALLINT NOT NULL CHECK (player_order IN (0, 1)),
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_room
        FOREIGN KEY(room_id)
        REFERENCES rooms(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_active_connection
        FOREIGN KEY(player_display_name)
        REFERENCES active_connections(display_name)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (room_id, player_order),
    UNIQUE (player_display_name)
);

CREATE INDEX idx_players_room_id ON players(room_id);
CREATE INDEX idx_players_player_display_name ON players(player_display_name);