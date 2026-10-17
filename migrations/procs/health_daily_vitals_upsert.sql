CREATE OR REPLACE PROCEDURE health_daily_vitals_upsert(
    p_user_id BIGINT,
    p_vital_date DATE,
    p_resting_heart_rate_bpm DOUBLE PRECISION,
    p_hrv_ms DOUBLE PRECISION,
    p_heart_rate_zone_1_minutes DOUBLE PRECISION,
    p_heart_rate_zone_2_minutes DOUBLE PRECISION,
    p_heart_rate_zone_3_minutes DOUBLE PRECISION,
    p_heart_rate_zone_4_minutes DOUBLE PRECISION,
    p_heart_rate_zone_5_minutes DOUBLE PRECISION
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_daily_vitals (
        user_id,
        vital_date,
        resting_heart_rate_bpm,
        hrv_ms,
        heart_rate_zone_1_minutes,
        heart_rate_zone_2_minutes,
        heart_rate_zone_3_minutes,
        heart_rate_zone_4_minutes,
        heart_rate_zone_5_minutes
    )
    VALUES (
        p_user_id,
        p_vital_date,
        p_resting_heart_rate_bpm,
        p_hrv_ms,
        p_heart_rate_zone_1_minutes,
        p_heart_rate_zone_2_minutes,
        p_heart_rate_zone_3_minutes,
        p_heart_rate_zone_4_minutes,
        p_heart_rate_zone_5_minutes
    )

    ON CONFLICT (user_id, vital_date)
    DO UPDATE SET
        resting_heart_rate_bpm = COALESCE(
            EXCLUDED.resting_heart_rate_bpm,
            health_daily_vitals.resting_heart_rate_bpm
        ),
        hrv_ms = COALESCE(
            EXCLUDED.hrv_ms,
            health_daily_vitals.hrv_ms
        ),
        heart_rate_zone_1_minutes = COALESCE(
            EXCLUDED.heart_rate_zone_1_minutes,
            health_daily_vitals.heart_rate_zone_1_minutes
        ),
        heart_rate_zone_2_minutes = COALESCE(
            EXCLUDED.heart_rate_zone_2_minutes,
            health_daily_vitals.heart_rate_zone_2_minutes
        ),
        heart_rate_zone_3_minutes = COALESCE(
            EXCLUDED.heart_rate_zone_3_minutes,
            health_daily_vitals.heart_rate_zone_3_minutes
        ),
        heart_rate_zone_4_minutes = COALESCE(
            EXCLUDED.heart_rate_zone_4_minutes,
            health_daily_vitals.heart_rate_zone_4_minutes
        ),
        heart_rate_zone_5_minutes = COALESCE(
            EXCLUDED.heart_rate_zone_5_minutes,
            health_daily_vitals.heart_rate_zone_5_minutes
        ),
        updated_at = CURRENT_TIMESTAMP;

END;
$$;