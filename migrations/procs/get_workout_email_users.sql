CREATE OR REPLACE FUNCTION get_workout_email_users()
RETURNS TABLE (
    user_id INT
)
LANGUAGE sql
AS $$
    SELECT user_id
    FROM workout_email_settings
    WHERE email <> ''
      AND email_time <> '';
$$;