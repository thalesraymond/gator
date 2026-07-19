-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    url VARCHAR NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE feeds;