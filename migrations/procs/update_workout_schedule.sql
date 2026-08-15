CREATE OR REPLACE PROCEDURE update_workout_schedule(
    p_workout_schedule_id INT,
    p_day_of_week SMALLINT,
    p_workout_name TEXT,
    p_exercises TEXT[]
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE workout_schedules
    SET
        day_of_week = p_day_of_week,
        workout_name = TRIM(p_workout_name),
        updated_at = CURRENT_TIMESTAMP
    WHERE workout_schedule_id = p_workout_schedule_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION
            'workout schedule with id % not found',
            p_workout_schedule_id
            USING ERRCODE = 'P0002';
    END IF;

    DELETE FROM workout_schedule_exercises
    WHERE workout_schedule_id = p_workout_schedule_id;

    INSERT INTO workout_schedule_exercises (
        workout_schedule_id,
        exercise_name,
        exercise_order
    )
    SELECT
        p_workout_schedule_id,
        TRIM(exercise_name),
        exercise_order::INT
    FROM UNNEST(COALESCE(p_exercises, ARRAY[]::TEXT[]))
        WITH ORDINALITY AS exercises(exercise_name, exercise_order)
    WHERE LENGTH(TRIM(exercise_name)) > 0;
END;
$$;
