CREATE TABLE health_workouts (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    google_data_point_id TEXT NOT NULL,

    workout_type TEXT,

    start_time TIMESTAMP NOT NULL,

    end_time TIMESTAMP NOT NULL,

    duration_minutes DOUBLE PRECISION,

    distance_meters DOUBLE PRECISION,

    active_energy_kcal DOUBLE PRECISION,

    total_calories_kcal DOUBLE PRECISION,

    average_heart_rate_bpm DOUBLE PRECISION,

    max_heart_rate_bpm DOUBLE PRECISION,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_workouts_time_check
        CHECK (end_time > start_time),

    CONSTRAINT health_workouts_google_id_unique
        UNIQUE (user_id, google_data_point_id)
);

CREATE INDEX idx_health_workouts_user_start_time
ON health_workouts(user_id, start_time);