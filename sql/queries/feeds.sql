-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, image, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.id AS feed_id, feeds.name, feeds.image, feeds.url, JSON_AGG((category.id, category.title)) AS categories 
FROM feeds
LEFT JOIN feed_categories ON feeds.id = feed_categories.feed_id
LEFT JOIN category ON feed_categories.category_id = category.id
GROUP BY feeds.id;
--SELECT * FROM feeds;

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
