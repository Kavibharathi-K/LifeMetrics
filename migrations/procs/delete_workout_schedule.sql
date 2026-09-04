CREATE OR REPLACE FUNCTION delete_workout_schedule(
    p_workout_schedule_id INT,
    p_user_id INT
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM workout_schedules
    WHERE workout_schedule_id = p_workout_schedule_id
      AND user_id = p_user_id;

    RETURN FOUND;
END;
$$;