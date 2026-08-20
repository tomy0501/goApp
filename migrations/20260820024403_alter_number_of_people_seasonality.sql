-- +goose Up
SELECT 'up SQL query';

ALTER TABLE `seasonality_number_of_people`
    MODIFY COLUMN `number_of_people` INT UNSIGNED NOT NULL DEFAULT 1;
-- +goose Down
SELECT 'down SQL query';

ALTER TABLE `seasonality_number_of_people`
    MODIFY COLUMN `number_of_people` INT NOT NULL;
