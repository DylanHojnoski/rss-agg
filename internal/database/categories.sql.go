// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: categories.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO category (id, created_at, updated_at, title)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at, title
`

type CreateCategoryParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, createCategory,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
	)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
	)
	return i, err
}

const getCategories = `-- name: GetCategories :many
SELECT id, created_at, updated_at, title FROM category
`

func (q *Queries) GetCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
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

const getCategoryByName = `-- name: GetCategoryByName :one
SELECT id, created_at, updated_at, title FROM category
WHERE title LIKE $1
`

func (q *Queries) GetCategoryByName(ctx context.Context, title string) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategoryByName, title)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
	)
	return i, err
}
