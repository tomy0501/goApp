-- +goose Up
SELECT 'up SQL query';


ALTER TABLE `room_type` ADD CONSTRAINT `room_type_facilities`  FOREIGN KEY (`facility_id`) REFERENCES facilities(id);

ALTER TABLE `room_type` 
MODIFY `facility_id` BINARY(16) NOT NULL;

ALTER TABLE `room` ADD CONSTRAINT `room_roomtype_id` FOREIGN KEY (`room_type_id`) REFERENCES room_type(id);

ALTER TABLE `room`
MODIFY `room_type_id` BINARY(16) NOT NULL;

-- +goose Down