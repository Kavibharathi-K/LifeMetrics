CREATE OR REPLACE FUNCTION health_sleep_session_upsert(
    p_user_id BIGINT,
    p_google_data_point_id TEXT,
    p_start_time TIMESTAMP,
    p_end_time TIMESTAMP,
    p_duration_minutes DOUBLE PRECISION,
    p_sleep_type TEXT,
    p_is_nap BOOLEAN,
    p_is_main_sleep BOOLEAN,
    p_processed BOOLEAN,
    p_manually_edited BOOLEAN,
    p_stages_status TEXT,
    p_external_id TEXT,
    p_minutes_in_sleep_period DOUBLE PRECISION,
    p_minutes_after_wake_up DOUBLE PRECISION,
    p_minutes_to_fall_asleep DOUBLE PRECISION,
    p_minutes_asleep DOUBLE PRECISION,
    p_minutes_awake DOUBLE PRECISION
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    v_sleep_session_id BIGINT;
BEGIN

    INSERT INTO health_sleep_sessions (
        user_id,
        google_data_point_id,
        start_time,
        end_time,
        duration_minutes,
        sleep_type,
        is_nap,
        is_main_sleep,
        processed,
        manually_edited,
        stages_status,
        external_id,
        minutes_in_sleep_period,
        minutes_after_wake_up,
        minutes_to_fall_asleep,
        minutes_asleep,
        minutes_awake
    )
    VALUES (
        p_user_id,
        p_google_data_point_id,
        p_start_time,
        p_end_time,
        p_duration_minutes,
        p_sleep_type,
        p_is_nap,
        p_is_main_sleep,
        p_processed,
        p_manually_edited,
        p_stages_status,
        p_external_id,
        p_minutes_in_sleep_period,
        p_minutes_after_wake_up,
        p_minutes_to_fall_asleep,
        p_minutes_asleep,
        p_minutes_awake
    )
    ON CONFLICT (user_id, google_data_point_id)
    DO UPDATE SET
        start_time = EXCLUDED.start_time,
        end_time = EXCLUDED.end_time,
        duration_minutes = EXCLUDED.duration_minutes,
        sleep_type = EXCLUDED.sleep_type,
        is_nap = EXCLUDED.is_nap,
        is_main_sleep = EXCLUDED.is_main_sleep,
        processed = EXCLUDED.processed,
        manually_edited = EXCLUDED.manually_edited,
        stages_status = EXCLUDED.stages_status,
        external_id = EXCLUDED.external_id,
        minutes_in_sleep_period = EXCLUDED.minutes_in_sleep_period,
        minutes_after_wake_up = EXCLUDED.minutes_after_wake_up,
        minutes_to_fall_asleep = EXCLUDED.minutes_to_fall_asleep,
        minutes_asleep = EXCLUDED.minutes_asleep,
        minutes_awake = EXCLUDED.minutes_awake,
        updated_at = CURRENT_TIMESTAMP

    RETURNING id INTO v_sleep_session_id;

    RETURN v_sleep_session_id;

END;
$$;