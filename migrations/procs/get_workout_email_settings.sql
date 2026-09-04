CREATE OR REPLACE FUNCTION get_workout_email_settings(
    p_user_id INT
)
RETURNS TABLE (
    email VARCHAR(254),
    email_time VARCHAR(5),
    updated_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        wes.email,
        wes.email_time,
        wes.updated_at
    FROM workout_email_settings wes
    WHERE wes.user_id = p_user_id;
$$;