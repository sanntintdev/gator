-- name: CreateFeedFollow :one
WITH inserted_follow AS (
    INSERT INTO feed_follows (user_id, feed_id, created_at, updated_at)
    VALUES ($1, $2, NOW(), NOW())
    RETURNING *
)
SELECT
    ff.id,
    ff.user_id,
    ff.feed_id,
    ff.created_at,
    ff.updated_at,
    u.name AS user_name,
    f.name AS feed_name
FROM inserted_follow ff
INNER JOIN users u ON ff.user_id = u.id
INNER JOIN feeds f ON ff.feed_id = f.id;

-- name: RetrieveFeedFollowsForUser :many
SELECT
    ff.id,
    ff.user_id,
    ff.feed_id,
    ff.created_at,
    ff.updated_at,
    u.name AS user_name,
    f.name AS feed_name
FROM feed_follows ff
INNER JOIN users u ON ff.user_id = u.id
INNER JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;
