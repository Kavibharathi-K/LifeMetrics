CREATE OR REPLACE FUNCTION update_workout_email_settings(
    p_user_id INT,
    p_email VARCHAR(254),
    p_email_time VARCHAR(5)
)
RETURNS TIMESTAMP
LANGUAGE plpgsql
AS $$
DECLARE
    v_updated_at TIMESTAMP;
BEGIN
    INSERT INTO workout_email_settings (
        user_id,
        email,
        email_time
    )
    VALUES (
        p_user_id,
        p_email,
        p_email_time
    )
    ON CONFLICT (user_id)
    DO UPDATE SET
        email = EXCLUDED.email,
        email_time = EXCLUDED.email_time,
        updated_at = CURRENT_TIMESTAMP
    RETURNING updated_at
    INTO v_updated_at;

    RETURN v_updated_at;
END;
$$;