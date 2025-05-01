-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, description, image, url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.id AS id, feeds.name, feeds.description, feeds.image, feeds.url, JSON_AGG((category.id, category.title)) AS categories 
FROM feeds
LEFT JOIN feed_categories ON feeds.id = feed_categories.feed_id
LEFT JOIN category ON feed_categories.category_id = category.id
GROUP BY feeds.id;

-- name: GetFeedForID :one
SELECT feeds.id AS id, feeds.name, feeds.description, feeds.image, feeds.url, JSON_AGG((category.id, category.title)) AS categories 
FROM feeds
LEFT JOIN feed_categories ON feeds.id = feed_categories.feed_id
LEFT JOIN category ON feed_categories.category_id = category.id
WHERE feeds.id = $1
GROUP BY feeds.id;

-- name: GetFeedsForCategory :many
SELECT feeds.id AS id, feeds.name, feeds.description, feeds.image, feeds.url, JSON_AGG((category.id, category.title)) AS categories 
FROM feeds
LEFT JOIN feed_categories ON feeds.id = feed_categories.feed_id
LEFT JOIN category ON feed_categories.category_id = category.id
WHERE category.id = $1
GROUP BY feeds.id;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetFollowedFeeds :many
SELECT feeds.id AS id, feeds.name, feeds.description, feeds.image, feeds.url, JSON_AGG((category.id, category.title)) AS categories 
FROM feeds
LEFT JOIN feed_follows ON feeds.id = feed_follows.feed_id
LEFT JOIN feed_categories ON feeds.id = feed_categories.feed_id
LEFT JOIN category ON feed_categories.category_id = category.id
WHERE feed_follows.user_id = $1
GROUP BY feeds.id;
