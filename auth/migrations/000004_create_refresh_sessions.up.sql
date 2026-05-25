CREATE TABLE refresh_sessions (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    refresh_token TEXT NOT NULL UNIQUE,

    user_agent TEXT,
    ip_address TEXT,

    expires_at TIMESTAMP NOT NULL,

    revoked BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);