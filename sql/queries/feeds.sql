-- name: CreateFeed :one
INSERT INTO feeds (url, name, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: RetrieveFeedsWithUser :many
SELECT * FROM feeds
LEFT JOIN users ON feeds.user_id = users.id;
