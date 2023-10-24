// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, audio, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at, title, description, published_at, audio, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description sql.NullString
	PublishedAt time.Time
	Audio       string
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Description,
		arg.PublishedAt,
		arg.Audio,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.PublishedAt,
		&i.Audio,
		&i.FeedID,
	)
	return i, err
}

const getPostsForFeed = `-- name: GetPostsForFeed :many
SELECT id, created_at, updated_at, title, description, published_at, audio, feed_id FROM posts 
WHERE feed_id = $1
ORDER BY published_at DESC 
LIMIT $2
`

type GetPostsForFeedParams struct {
	FeedID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForFeed(ctx context.Context, arg GetPostsForFeedParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForFeed, arg.FeedID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.PublishedAt,
			&i.Audio,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsForFeedBeforeDate = `-- name: GetPostsForFeedBeforeDate :many
SELECT id, created_at, updated_at, title, description, published_at, audio, feed_id FROM posts
WHERE feed_id = $1 AND published_at < $2
ORDER BY published_at DESC 
LIMIT $3
`

type GetPostsForFeedBeforeDateParams struct {
	FeedID      uuid.UUID
	PublishedAt time.Time
	Limit       int32
}

func (q *Queries) GetPostsForFeedBeforeDate(ctx context.Context, arg GetPostsForFeedBeforeDateParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForFeedBeforeDate, arg.FeedID, arg.PublishedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.PublishedAt,
			&i.Audio,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.description, posts.published_at, posts.audio, posts.feed_id FROM posts 
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC 
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.PublishedAt,
			&i.Audio,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
