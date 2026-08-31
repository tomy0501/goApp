-- +goose Up
SELECT 'up SQL query';


ALTER TABLE `room` DROP FOREIGN KEY `fk_1`;

ALTER TABLE `room`
DROP COLUMN `facility_id`;

ALTER TABLE `room`
ADD COLUMN `room_type_id` BINARY(16);

-- +goose Down