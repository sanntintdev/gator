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

-- name: RetrieveFeedWithURL :one
SELECT * FROM feeds
WHERE  url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: RetrieveNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
