-- +goose Up
SELECT 'up SQL query';

ALTER TABLE `seasonality_of_date`
    ADD COLUMN `room_type_id` binary(16) NOT NULL;

ALTER TABLE `seasonality_of_date`
    ADD FOREIGN KEY (`room_type_id`) REFERENCES `room_type` (`id`);

ALTER TABLE `seasonality_of_date`
    ADD COLUMN `default_price` decimal(10, 2);

-- +goose Down
SELECT 'down SQL query';
