CREATE OR REPLACE FUNCTION create_workout_schedule(
    p_user_id INT,
    p_day_of_week SMALLINT,
    p_workout_name TEXT,
    p_exercises TEXT[]
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    v_workout_schedule_id INT;
BEGIN
    INSERT INTO workout_schedules (
        user_id,
        day_of_week,
        workout_name
    )
    VALUES (
        p_user_id,
        p_day_of_week,
        TRIM(p_workout_name)
    )
    RETURNING workout_schedule_id
    INTO v_workout_schedule_id;

    INSERT INTO workout_schedule_exercises (
        workout_schedule_id,
        exercise_name,
        exercise_order
    )
    SELECT
        v_workout_schedule_id,
        TRIM(exercise_name),
        exercise_order::INT
    FROM UNNEST(COALESCE(p_exercises, ARRAY[]::TEXT[]))
        WITH ORDINALITY AS exercises(exercise_name, exercise_order)
    WHERE LENGTH(TRIM(exercise_name)) > 0;

    RETURN v_workout_schedule_id;
END;
$$;