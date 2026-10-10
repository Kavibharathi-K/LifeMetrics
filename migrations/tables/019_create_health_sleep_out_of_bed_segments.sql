CREATE TABLE health_sleep_out_of_bed_segments (
    id BIGSERIAL PRIMARY KEY,

    sleep_session_id BIGINT NOT NULL
        REFERENCES health_sleep_sessions(id)
        ON DELETE CASCADE,

    start_time TIMESTAMP NOT NULL,

    end_time TIMESTAMP NOT NULL,

    duration_minutes DOUBLE PRECISION,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_sleep_out_of_bed_time_check
        CHECK (end_time > start_time)
);

CREATE INDEX idx_health_sleep_out_of_bed_session_id
ON health_sleep_out_of_bed_segments(sleep_session_id);