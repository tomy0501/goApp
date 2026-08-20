-- +goose Up
SELECT 'up SQL query';

DROP TABLE `seasonality_price`;

ALTER TABLE `seasonality_number_of_people` 
    ADD COLUMN `price` decimal(10, 2) NOT NULL DEFAULT 0.00;
-- +goose Down
SELECT 'down SQL query';
CREATE TABLE `seasonality_price` (
    `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(),1)),
    `seasonality_number_of_people_id` binary(16) NOT NULL,
    `price` decimal(10, 2) NOT NULL,
    `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
    FOREIGN KEY (`seasonality_number_of_people_id`) REFERENCES `seasonality_number_of_people` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
