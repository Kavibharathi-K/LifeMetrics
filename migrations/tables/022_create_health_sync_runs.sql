CREATE TABLE health_sync_runs (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    provider TEXT NOT NULL,

    sync_type TEXT NOT NULL,

    started_at TIMESTAMP NOT NULL,

    completed_at TIMESTAMP,

    status TEXT NOT NULL,

    error_message TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_sync_runs_status_check
        CHECK (status IN ('RUNNING', 'SUCCESS', 'FAILED'))
);

CREATE INDEX idx_health_sync_runs_user_started_at
ON health_sync_runs(user_id, started_at DESC);