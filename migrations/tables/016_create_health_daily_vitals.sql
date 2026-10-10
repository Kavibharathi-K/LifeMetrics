CREATE TABLE health_daily_vitals (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    vital_date DATE NOT NULL,

    resting_heart_rate_bpm DOUBLE PRECISION,

    hrv_ms DOUBLE PRECISION,

    heart_rate_zone_1_minutes DOUBLE PRECISION,
    heart_rate_zone_2_minutes DOUBLE PRECISION,
    heart_rate_zone_3_minutes DOUBLE PRECISION,
    heart_rate_zone_4_minutes DOUBLE PRECISION,
    heart_rate_zone_5_minutes DOUBLE PRECISION,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT health_daily_vitals_user_date_unique
        UNIQUE (user_id, vital_date)
);