-- CREATE DATABASE subscribers;

-- +goose Up
CREATE TABLE usr_subscription(
    id SERIAL PRIMARY KEY,
    service text NOT NULL,
    usr_id uuid NOT NULL,
    price INTEGER check ( >=0 ) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP
);

CREATE INDEX usr_subscription_idx ON usr_subscription (usr_id, start_date, end_date);

-- +goose Down
DROP INDEX usr_subscription_idx;
DROP TABLE usr_subscription;