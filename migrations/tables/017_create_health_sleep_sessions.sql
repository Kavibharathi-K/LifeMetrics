CREATE TABLE health_sleep_sessions (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    google_data_point_id TEXT NOT NULL,

    start_time TIMESTAMP NOT NULL,

    end_time TIMESTAMP NOT NULL,

    duration_minutes DOUBLE PRECISION,

    sleep_type TEXT,

    is_nap BOOLEAN,

    is_main_sleep BOOLEAN,

    processed BOOLEAN,

    manually_edited BOOLEAN,

    stages_status TEXT,

    external_id TEXT,

    minutes_in_sleep_period DOUBLE PRECISION,

    minutes_after_wake_up DOUBLE PRECISION,

    minutes_to_fall_asleep DOUBLE PRECISION,

    minutes_asleep DOUBLE PRECISION,

    minutes_awake DOUBLE PRECISION,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_sleep_sessions_time_check
        CHECK (end_time > start_time),

    CONSTRAINT health_sleep_sessions_google_id_unique
        UNIQUE (user_id, google_data_point_id)
);