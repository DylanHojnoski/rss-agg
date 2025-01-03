-- name: CreateCategory :one
INSERT INTO category (id, created_at, updated_at, title)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCategories :many
SELECT * FROM category;

-- name: GetCategoryByName :one
SELECT * FROM category
WHERE title LIKE $1;
