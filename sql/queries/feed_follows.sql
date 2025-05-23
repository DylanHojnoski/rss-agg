-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows
WHERE user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows 
WHERE feed_id = $1 AND user_id = $2;

-- name: CreateFeedFollowWithURL :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, 
    (SELECT id FROM feeds WHERE feeds.url = $5)
)
RETURNING *;
