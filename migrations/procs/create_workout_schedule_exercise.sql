DROP FUNCTION IF EXISTS create_workout_schedule_exercise(INT, TEXT, INT);

CREATE OR REPLACE FUNCTION create_workout_schedule_exercise(
    p_user_id INT,
    p_workout_schedule_id INT,
    p_exercise_name TEXT,
    p_exercise_order INT DEFAULT NULL
)
RETURNS TABLE (
    workout_schedule_exercise_id INT,
    exercise_order INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_exercise_order INT;
BEGIN
    /*
     * Verify that the workout schedule exists
     * and belongs to the authenticated user.
     */
    PERFORM 1
    FROM workout_schedules
    WHERE workout_schedule_id = p_workout_schedule_id
      AND user_id = p_user_id
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION
            'workout schedule with id % not found for user',
            p_workout_schedule_id
            USING ERRCODE = 'P0002';
    END IF;

    /*
     * If exercise_order is not provided,
     * append the exercise to the end.
     */
    IF p_exercise_order IS NULL OR p_exercise_order = 0 THEN

        SELECT COALESCE(MAX(wse.exercise_order), 0) + 1
        INTO v_exercise_order
        FROM workout_schedule_exercises wse
        WHERE wse.workout_schedule_id = p_workout_schedule_id;

    ELSE

        v_exercise_order := p_exercise_order;

    END IF;

    RETURN QUERY
    INSERT INTO workout_schedule_exercises (
        workout_schedule_id,
        exercise_name,
        exercise_order
    )
    VALUES (
        p_workout_schedule_id,
        TRIM(p_exercise_name),
        v_exercise_order
    )
    RETURNING
        workout_schedule_exercises.workout_schedule_exercise_id,
        workout_schedule_exercises.exercise_order;
END;
$$;