CREATE OR REPLACE FUNCTION health_daily_activity_get_by_date(
    p_user_id BIGINT,
    p_activity_date DATE
)
RETURNS TABLE (
    id BIGINT,
    activity_date DATE,
    steps BIGINT,
    distance_meters DOUBLE PRECISION,
    active_energy_kcal DOUBLE PRECISION,
    total_calories_kcal DOUBLE PRECISION,
    active_minutes DOUBLE PRECISION,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        h.id,
        h.activity_date,
        h.steps,
        h.distance_meters,
        h.active_energy_kcal,
        h.total_calories_kcal,
        h.active_minutes,
        h.created_at,
        h.updated_at
    FROM health_daily_activity h
    WHERE h.user_id = p_user_id
      AND h.activity_date = p_activity_date;
$$;