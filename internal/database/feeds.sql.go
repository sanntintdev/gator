// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (url, name, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, url, name, user_id, created_at, updated_at, last_fetched_at
`

type CreateFeedParams struct {
	Url       string
	Name      string
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.Url,
		arg.Name,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1
`

func (q *Queries) MarkFeedFetched(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, id)
	return err
}

const retrieveFeedWithURL = `-- name: RetrieveFeedWithURL :one
SELECT id, url, name, user_id, created_at, updated_at, last_fetched_at FROM feeds
WHERE  url = $1
`

func (q *Queries) RetrieveFeedWithURL(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, retrieveFeedWithURL, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}

const retrieveFeedsWithUser = `-- name: RetrieveFeedsWithUser :many
SELECT feeds.id, url, feeds.name, user_id, feeds.created_at, feeds.updated_at, last_fetched_at, users.id, users.name, users.created_at, users.updated_at FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
`

type RetrieveFeedsWithUserRow struct {
	ID            int32
	Url           string
	Name          string
	UserID        uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt sql.NullTime
	ID_2          uuid.NullUUID
	Name_2        sql.NullString
	CreatedAt_2   sql.NullTime
	UpdatedAt_2   sql.NullTime
}

func (q *Queries) RetrieveFeedsWithUser(ctx context.Context) ([]RetrieveFeedsWithUserRow, error) {
	rows, err := q.db.QueryContext(ctx, retrieveFeedsWithUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RetrieveFeedsWithUserRow
	for rows.Next() {
		var i RetrieveFeedsWithUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Name,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastFetchedAt,
			&i.ID_2,
			&i.Name_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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

const retrieveNextFeedToFetch = `-- name: RetrieveNextFeedToFetch :one
SELECT id, url, name, user_id, created_at, updated_at, last_fetched_at FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1
`

func (q *Queries) RetrieveNextFeedToFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, retrieveNextFeedToFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastFetchedAt,
	)
	return i, err
}
