ALTER TABLE rooms
ADD COLUMN host_display_name VARCHAR(50) NOT NULL;

ALTER TABLE rooms
ADD CONSTRAINT fk_host_connection
    FOREIGN KEY(host_display_name)
    REFERENCES active_connections(display_name)
    ON DELETE CASCADE;
