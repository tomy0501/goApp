-- +goose Up
SELECT 'up SQL query';

ALTER TABLE room_type
ADD COLUMN seasonality_id BINARY(16);

ALTER TABLE room_type ADD CONSTRAINT fk_room_type_seasonality FOREIGN KEY (seasonality_id) REFERENCES seasonality(id);
-- +goose Down
SELECT 'down SQL query';
