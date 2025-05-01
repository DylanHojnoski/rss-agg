-- +goose Up
CREATE TABLE category (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE category;
