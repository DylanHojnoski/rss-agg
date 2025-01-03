-- name: CreateFeedCategory :one
INSERT INTO feed_categories (id, created_at, updated_at, feed_id, category_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedsForCategory :many
SELECT * FROM feed_categories
JOIN feeds ON feeds.id = feed_categories.feed_id
WHERE category_id = $1;

