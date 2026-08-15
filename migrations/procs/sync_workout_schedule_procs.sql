DROP FUNCTION IF EXISTS create_workout_schedule(SMALLINT, TEXT, TEXT, TEXT, TEXT[]);
DROP PROCEDURE IF EXISTS update_workout_schedule(INT, SMALLINT, TEXT, TEXT, TEXT, TEXT[]);

CREATE OR REPLACE FUNCTION create_workout_schedule(
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
        day_of_week,
        workout_name
    )
    VALUES (
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

CREATE OR REPLACE FUNCTION list_workout_schedules()
RETURNS TABLE (
    workout_schedule_id INT,
    day_of_week SMALLINT,
    workout_name VARCHAR(100),
    exercises JSONB,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        ws.workout_schedule_id,
        ws.day_of_week,
        ws.workout_name,
        COALESCE(
            JSONB_AGG(
                JSONB_BUILD_OBJECT(
                    'workout_schedule_exercise_id',
                    wse.workout_schedule_exercise_id,
                    'exercise_name',
                    wse.exercise_name,
                    'exercise_order',
                    wse.exercise_order
                )
                ORDER BY wse.exercise_order
            ) FILTER (
                WHERE wse.workout_schedule_exercise_id IS NOT NULL
            ),
            '[]'::JSONB
        ) AS exercises,
        ws.created_at,
        ws.updated_at
    FROM workout_schedules ws
    LEFT JOIN workout_schedule_exercises wse
        ON wse.workout_schedule_id = ws.workout_schedule_id
    GROUP BY
        ws.workout_schedule_id,
        ws.day_of_week,
        ws.workout_name,
        ws.created_at,
        ws.updated_at
    ORDER BY ws.day_of_week;
$$;

CREATE OR REPLACE FUNCTION get_workout_schedule_by_day(
    p_day_of_week SMALLINT
)
RETURNS TABLE (
    workout_schedule_id INT,
    day_of_week SMALLINT,
    workout_name VARCHAR(100),
    exercises JSONB,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        ws.workout_schedule_id,
        ws.day_of_week,
        ws.workout_name,
        COALESCE(
            JSONB_AGG(
                JSONB_BUILD_OBJECT(
                    'workout_schedule_exercise_id',
                    wse.workout_schedule_exercise_id,
                    'exercise_name',
                    wse.exercise_name,
                    'exercise_order',
                    wse.exercise_order
                )
                ORDER BY wse.exercise_order
            ) FILTER (
                WHERE wse.workout_schedule_exercise_id IS NOT NULL
            ),
            '[]'::JSONB
        ) AS exercises,
        ws.created_at,
        ws.updated_at
    FROM workout_schedules ws
    LEFT JOIN workout_schedule_exercises wse
        ON wse.workout_schedule_id = ws.workout_schedule_id
    WHERE ws.day_of_week = p_day_of_week
    GROUP BY
        ws.workout_schedule_id,
        ws.day_of_week,
        ws.workout_name,
        ws.created_at,
        ws.updated_at;
$$;
