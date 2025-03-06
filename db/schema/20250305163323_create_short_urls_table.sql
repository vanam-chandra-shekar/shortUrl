-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shorturls (
    sid SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shorturls
-- +goose StatementEnd
