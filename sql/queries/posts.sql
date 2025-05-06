-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, audio, duration, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name,
CASE 
    WHEN post_views.id IS NOT NULL THEN TRUE
    ELSE FALSE
END AS viewed
FROM posts 
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
JOIN feeds ON posts.feed_id = feeds.id
LEFT JOIN post_views ON posts.id = post_views.post_id
WHERE feed_follows.user_id = $1 AND post_views.id IS NULL
ORDER BY posts.published_at DESC 
LIMIT $2;

-- name: GetPostsForFeed :many 
SELECT posts.*, feeds.name AS feed_name,
CASE 
    WHEN post_views.id IS NOT NULL THEN TRUE
    ELSE FALSE
END AS viewed
FROM posts
JOIN feeds ON posts.feed_id = feeds.id
LEFT JOIN post_views ON posts.id = post_views.post_id AND (sqlc.arg(UserID)::uuid IS NULL OR sqlc.arg(UserID)::uuid = post_views.user_id)
WHERE feed_id = $1 AND (sqlc.arg(Unviewed)::bool = FALSE OR post_views.id IS NULL)
ORDER BY published_at DESC 
LIMIT $2;

-- name: GetPostsForFeedBeforeDate :many 
SELECT posts.*, feeds.name AS feed_name,
CASE 
    WHEN post_views.id IS NOT NULL THEN TRUE
    ELSE FALSE
END AS viewed
FROM posts
JOIN feeds ON posts.feed_id = feeds.id
LEFT JOIN post_views ON posts.id = post_views.post_id AND (sqlc.arg(UserID)::uuid IS NULL OR sqlc.arg(UserID)::uuid = post_views.user_id)
WHERE feed_id = $1 AND published_at < $2 AND (sqlc.arg(Unviewed)::bool = FALSE OR post_views.id IS NULL)
ORDER BY published_at DESC 
LIMIT $3;

-- name: GetPostsForFeedAsc :many 
SELECT posts.*, feeds.name AS feed_name,
CASE 
    WHEN post_views.id IS NOT NULL THEN TRUE
    ELSE FALSE
END AS viewed
FROM posts 
JOIN feeds ON posts.feed_id = feeds.id
LEFT JOIN post_views ON posts.id = post_views.post_id AND (sqlc.arg(UserID)::uuid IS NULL OR sqlc.arg(UserID)::uuid = post_views.user_id)
WHERE feed_id = $1 AND (sqlc.arg(Unviewed)::bool = FALSE OR post_views.id IS NULL)
ORDER BY published_at Asc 
LIMIT $2;

-- name: GetPostsForFeedAfterDate :many 
SELECT posts.*, feeds.name AS feed_name,
CASE 
    WHEN post_views.id IS NOT NULL THEN TRUE
    ELSE FALSE
END AS viewed
FROM posts
JOIN feeds ON posts.feed_id = feeds.id
LEFT JOIN post_views ON posts.id = post_views.post_id AND (sqlc.arg(UserID)::uuid IS NULL OR sqlc.arg(UserID)::uuid = post_views.user_id)
WHERE feed_id = $1 AND published_at > $2 AND (sqlc.arg(Unviewed)::bool = FALSE OR post_views.id IS NULL)
ORDER BY published_at ASC 
LIMIT $3;


