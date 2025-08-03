-- +goose Up
CREATE TABLE auth_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES auth_sessions(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL CHECK (event_type IN ('login', 'logout', 'refresh', 'revoked')),
    ip_address TEXT,
    user_agent TEXT,
    event_time TIMESTAMP NOT NULL DEFAULT now()
);

-- Create indexes for performance and analytics
CREATE INDEX idx_auth_logs_user_id ON auth_logs(user_id);
CREATE INDEX idx_auth_logs_session_id ON auth_logs(session_id);
CREATE INDEX idx_auth_logs_event_type ON auth_logs(event_type);
CREATE INDEX idx_auth_logs_event_time ON auth_logs(event_time);

-- +goose Down
DROP TABLE IF EXISTS auth_logs;
