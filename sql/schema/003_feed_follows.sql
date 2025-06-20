-- +goose Up
CREATE TABLE feed_follows (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id),
    feed_id SERIAL NOT NULL REFERENCES feeds (id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
