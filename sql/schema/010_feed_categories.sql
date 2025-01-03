-- +goose Up
CREATE TABLE feed_categories (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    CONSTRAINT unique_category_pair UNIQUE (feed_id, category_id)
);

-- +goose Down
DROP TABLE feed_categories;
