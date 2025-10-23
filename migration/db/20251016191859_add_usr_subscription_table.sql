-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS usr_subscription(
                                               id SERIAL PRIMARY KEY,
                                               service text NOT NULL,
                                               usr_id uuid NOT NULL,
                                               price INTEGER NOT NULL,
                                               start_date TIMESTAMP NOT NULL,
                                               end_date TIMESTAMP
);
CREATE INDEX IF NOT EXISTS usr_subscription_idx ON usr_subscription (usr_id, start_date, end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX usr_subscription_idx;
DROP TABLE usr_subscription;
-- +goose StatementEnd
