CREATE OR REPLACE PROCEDURE health_daily_activity_upsert(
    p_user_id BIGINT,
    p_activity_date DATE,
    p_steps BIGINT,
    p_distance_meters DOUBLE PRECISION,
    p_active_energy_kcal DOUBLE PRECISION,
    p_total_calories_kcal DOUBLE PRECISION,
    p_active_minutes DOUBLE PRECISION
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_daily_activity (
        user_id,
        activity_date,
        steps,
        distance_meters,
        active_energy_kcal,
        total_calories_kcal,
        active_minutes
    )
    VALUES (
        p_user_id,
        p_activity_date,
        p_steps,
        p_distance_meters,
        p_active_energy_kcal,
        p_total_calories_kcal,
        p_active_minutes
    )

    ON CONFLICT (user_id, activity_date)
    DO UPDATE SET
        steps = EXCLUDED.steps,
        distance_meters = EXCLUDED.distance_meters,
        active_energy_kcal = EXCLUDED.active_energy_kcal,
        total_calories_kcal = EXCLUDED.total_calories_kcal,
        active_minutes = EXCLUDED.active_minutes,
        updated_at = CURRENT_TIMESTAMP;

END;
$$;