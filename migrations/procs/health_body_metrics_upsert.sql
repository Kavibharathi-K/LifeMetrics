CREATE OR REPLACE PROCEDURE health_body_metric_upsert(
    p_user_id BIGINT,
    p_google_data_point_id TEXT,
    p_metric_type TEXT,
    p_value DOUBLE PRECISION,
    p_unit TEXT,
    p_measured_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_body_metrics (
        user_id,
        google_data_point_id,
        metric_type,
        value,
        unit,
        measured_at
    )
    VALUES (
        p_user_id,
        p_google_data_point_id,
        p_metric_type,
        p_value,
        p_unit,
        p_measured_at
    )
    ON CONFLICT (user_id, google_data_point_id)
    DO UPDATE SET
        metric_type = EXCLUDED.metric_type,
        value = EXCLUDED.value,
        unit = EXCLUDED.unit,
        measured_at = EXCLUDED.measured_at,
        updated_at = CURRENT_TIMESTAMP;

END;
$$;