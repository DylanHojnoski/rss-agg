-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    url TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE feeds;
