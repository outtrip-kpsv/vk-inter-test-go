-- +goose Up
ALTER TABLE users
    ALTER COLUMN role_id SET DEFAULT 1;

-- +goose Down
ALTER TABLE users
    ALTER COLUMN role_id DROP DEFAULT;