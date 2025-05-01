-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL
);

-- +goose Down
DROP TABLE users;
