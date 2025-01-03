// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
}

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	LastFetchedAt sql.NullTime
	Image         sql.NullString
}

type FeedCategory struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	FeedID     uuid.UUID
	CategoryID uuid.UUID
}

type FeedFollow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type Listened struct {
	ID     uuid.UUID
	UserID uuid.UUID
	FeedID uuid.UUID
}

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description sql.NullString
	PublishedAt time.Time
	Audio       string
	Duration    sql.NullString
	FeedID      uuid.UUID
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ApiKey    string
}
