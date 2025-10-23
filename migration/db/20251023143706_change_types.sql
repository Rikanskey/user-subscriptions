-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE usr_subscription ALTER COLUMN start_date TYPE TIMESTAMP WITHOUT TIME ZONE;
ALTER TABLE usr_subscription ALTER COLUMN end_date TYPE TIMESTAMP WITHOUT TIME ZONE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE usr_subscription ALTER COLUMN start_date TYPE TIMESTAMP;
ALTER TABLE usr_subscription ALTER COLUMN end_date TYPE TIMESTAMP;
-- +goose StatementEnd
