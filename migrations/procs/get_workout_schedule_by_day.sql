CREATE OR REPLACE FUNCTION get_workout_schedule_by_day(
    p_user_id INT,
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
    WHERE ws.user_id = p_user_id
      AND ws.day_of_week = p_day_of_week
    GROUP BY
        ws.workout_schedule_id,
        ws.day_of_week,
        ws.workout_name,
        ws.created_at,
        ws.updated_at;
$$;