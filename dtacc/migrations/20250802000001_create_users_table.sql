-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Create index for email lookups
CREATE INDEX idx_users_email ON users(email);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "uuid-ossp";
