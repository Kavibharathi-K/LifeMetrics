CREATE OR REPLACE FUNCTION user_add(
    p_email TEXT,
    p_password_hash TEXT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    v_id BIGINT;
BEGIN
    INSERT INTO users (
        email,
        password_hash
    )
    VALUES (
        p_email,
        p_password_hash
    )
    RETURNING id
    INTO v_id;

    RETURN v_id;
END;
$$;