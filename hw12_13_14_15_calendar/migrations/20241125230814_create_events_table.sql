-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE events (
                        id UUID PRIMARY KEY,
                        title TEXT NOT NULL,
                        event_date TIMESTAMP NOT NULL,
                        duration INTERVAL NOT NULL,
                        description TEXT,
                        user_id INT NOT NULL,
                        notify_before INTERVAL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
