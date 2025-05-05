-- name: CreatePostView :one
INSERT INTO post_views(id, created_at, updated_at, user_id, post_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
