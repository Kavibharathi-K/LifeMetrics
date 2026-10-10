CREATE OR REPLACE PROCEDURE health_workout_upsert(
    p_user_id BIGINT,
    p_google_data_point_id TEXT,
    p_workout_type TEXT,
    p_start_time TIMESTAMP,
    p_end_time TIMESTAMP,
    p_duration_minutes DOUBLE PRECISION,
    p_distance_meters DOUBLE PRECISION,
    p_active_energy_kcal DOUBLE PRECISION,
    p_total_calories_kcal DOUBLE PRECISION,
    p_average_heart_rate_bpm DOUBLE PRECISION,
    p_max_heart_rate_bpm DOUBLE PRECISION
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_workouts (
        user_id,
        google_data_point_id,
        workout_type,
        start_time,
        end_time,
        duration_minutes,
        distance_meters,
        active_energy_kcal,
        total_calories_kcal,
        average_heart_rate_bpm,
        max_heart_rate_bpm
    )
    VALUES (
        p_user_id,
        p_google_data_point_id,
        p_workout_type,
        p_start_time,
        p_end_time,
        p_duration_minutes,
        p_distance_meters,
        p_active_energy_kcal,
        p_total_calories_kcal,
        p_average_heart_rate_bpm,
        p_max_heart_rate_bpm
    )
    ON CONFLICT (user_id, google_data_point_id)
    DO UPDATE SET
        workout_type = EXCLUDED.workout_type,
        start_time = EXCLUDED.start_time,
        end_time = EXCLUDED.end_time,
        duration_minutes = EXCLUDED.duration_minutes,
        distance_meters = EXCLUDED.distance_meters,
        active_energy_kcal = EXCLUDED.active_energy_kcal,
        total_calories_kcal = EXCLUDED.total_calories_kcal,
        average_heart_rate_bpm = EXCLUDED.average_heart_rate_bpm,
        max_heart_rate_bpm = EXCLUDED.max_heart_rate_bpm,
        updated_at = CURRENT_TIMESTAMP;

END;
$$;