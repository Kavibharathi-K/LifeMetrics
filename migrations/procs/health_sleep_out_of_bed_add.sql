CREATE OR REPLACE PROCEDURE health_sleep_out_of_bed_segment_add(
    p_sleep_session_id BIGINT,
    p_start_time TIMESTAMP,
    p_end_time TIMESTAMP,
    p_duration_minutes DOUBLE PRECISION
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_sleep_out_of_bed_segments (
        sleep_session_id,
        start_time,
        end_time,
        duration_minutes
    )
    VALUES (
        p_sleep_session_id,
        p_start_time,
        p_end_time,
        p_duration_minutes
    );

END;
$$;