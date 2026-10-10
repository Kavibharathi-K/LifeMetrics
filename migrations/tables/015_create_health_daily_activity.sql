CREATE TABLE health_daily_activity (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    activity_date DATE NOT NULL,

    steps BIGINT,

    distance_meters DOUBLE PRECISION,

    active_energy_kcal DOUBLE PRECISION,

    total_calories_kcal DOUBLE PRECISION,

    active_minutes DOUBLE PRECISION,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_daily_activity_user_date_unique
        UNIQUE (user_id, activity_date)
);