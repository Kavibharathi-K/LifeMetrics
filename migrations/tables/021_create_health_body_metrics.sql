CREATE TABLE health_body_metrics (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    google_data_point_id TEXT NOT NULL,

    metric_type TEXT NOT NULL,

    value DOUBLE PRECISION NOT NULL,

    unit TEXT NOT NULL,

    measured_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_body_metrics_google_id_unique
        UNIQUE (user_id, google_data_point_id)
);

CREATE INDEX idx_health_body_metrics_user_type_time
ON health_body_metrics(user_id, metric_type, measured_at);