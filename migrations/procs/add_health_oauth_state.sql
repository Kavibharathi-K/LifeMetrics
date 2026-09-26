CREATE OR REPLACE PROCEDURE health_oauth_state_add(
    p_state TEXT,
    p_user_id BIGINT,
    p_expires_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_oauth_states (
        state,
        user_id,
        expires_at
    )
    VALUES (
        p_state,
        p_user_id,
        p_expires_at
    );

END;
$$;