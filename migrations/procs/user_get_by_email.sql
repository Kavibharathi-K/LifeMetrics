CREATE OR REPLACE FUNCTION user_get_by_email(
    p_email TEXT
)
RETURNS TABLE (
    id BIGINT,
    email TEXT,
    password_hash TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        u.id,
        u.email,
        u.password_hash,
        u.created_at,
        u.updated_at
    FROM users u
    WHERE u.email = p_email;
END;
$$;