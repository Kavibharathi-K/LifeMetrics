CREATE OR REPLACE FUNCTION health_sleep_stage_upsert(
    p_sleep_session_id BIGINT,
    p_stage TEXT,
    p_start_time TIMESTAMP,
    p_end_time TIMESTAMP,
    p_duration_minutes DOUBLE PRECISION
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    v_sleep_stage_id BIGINT;
BEGIN
    INSERT INTO health_sleep_stages (
        sleep_session_id,
        stage,
        start_time,
        end_time,
        duration_minutes
    )
    VALUES (
        p_sleep_session_id,
        p_stage,
        p_start_time,
        p_end_time,
        p_duration_minutes
    )
    ON CONFLICT (
        sleep_session_id,
        start_time,
        end_time
    )
    DO UPDATE SET
        stage = EXCLUDED.stage,
        duration_minutes = EXCLUDED.duration_minutes
    RETURNING id INTO v_sleep_stage_id;

    RETURN v_sleep_stage_id;
END;
$$;