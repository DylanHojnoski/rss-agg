-- +goose Up 
ALTER TABLE feeds ADD COLUMN image TEXT UNIQUE;

-- +goose Down
ALTER TABLE feeds DROP COLUMN image;
