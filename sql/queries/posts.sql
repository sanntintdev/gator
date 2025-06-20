-- name: CreatePost :one
INSERT INTO posts (title, url, description, published_at, feed_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, Now(), Now())
RETURNING id;

-- name: RetrievePostsForUser :many
SELECT id, title, url, description, published_at, feed_id, created_at, updated_at
FROM posts
ORDER BY published_at DESC
LIMIT $1;
