-- +goose Up
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_link TEXT NOT NULL,
    code TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE links;
