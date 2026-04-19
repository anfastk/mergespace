-- Create outbox_events table

CREATE TABLE IF NOT EXISTS outbox_events (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMP DEFAULT NOW(),
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for fast polling
CREATE INDEX IF NOT EXISTS idx_outbox_status_created ON outbox_events (
    status,
    next_retry_at,
    created_at
);

-- Optional: index for retry handling
CREATE INDEX IF NOT EXISTS idx_outbox_retry ON outbox_events (retry_count);