-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, username, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1
LIMIT 1;
