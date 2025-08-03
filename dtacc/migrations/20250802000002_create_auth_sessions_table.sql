-- +goose Up
CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash TEXT NOT NULL,
    user_agent TEXT,
    ip_address TEXT,
    device_name TEXT,
    is_valid BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX idx_auth_sessions_refresh_token_hash ON auth_sessions(refresh_token_hash);
CREATE INDEX idx_auth_sessions_is_valid ON auth_sessions(is_valid);
CREATE INDEX idx_auth_sessions_expires_at ON auth_sessions(expires_at);

-- +goose Down
DROP TABLE IF EXISTS auth_sessions;
