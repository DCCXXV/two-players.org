ALTER TABLE rooms
DROP CONSTRAINT IF EXISTS fk_host_connection;

ALTER TABLE rooms
DROP COLUMN IF EXISTS host_display_name;
