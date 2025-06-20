-- +goose Up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    published_at TIMESTAMP,
    feed_id SERIAL NOT NULL REFERENCES feeds(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE posts;
