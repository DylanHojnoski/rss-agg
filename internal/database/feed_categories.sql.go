// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feed_categories.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedCategory = `-- name: CreateFeedCategory :one
INSERT INTO feed_categories (id, created_at, updated_at, feed_id, category_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, feed_id, category_id
`

type CreateFeedCategoryParams struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	FeedID     uuid.UUID
	CategoryID uuid.UUID
}

func (q *Queries) CreateFeedCategory(ctx context.Context, arg CreateFeedCategoryParams) (FeedCategory, error) {
	row := q.db.QueryRowContext(ctx, createFeedCategory,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FeedID,
		arg.CategoryID,
	)
	var i FeedCategory
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.CategoryID,
	)
	return i, err
}
