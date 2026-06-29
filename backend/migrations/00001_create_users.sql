-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_email_not_empty CHECK (btrim(email) <> ''),
    CONSTRAINT users_password_hash_not_empty CHECK (btrim(password_hash) <> '')
);

CREATE UNIQUE INDEX users_email_lower_unique_idx ON users (lower(email));

-- +goose Down
DROP INDEX IF EXISTS users_email_lower_unique_idx;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS pgcrypto;
