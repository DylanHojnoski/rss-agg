-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, audio, duration, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts 
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC 
LIMIT $2;

-- name: GetPostsForFeed :many 
SELECT * FROM posts 
WHERE feed_id = $1
ORDER BY published_at DESC 
LIMIT $2;

-- name: GetPostsForFeedBeforeDate :many 
SELECT * FROM posts
WHERE feed_id = $1 AND published_at < $2
ORDER BY published_at DESC 
LIMIT $3;
