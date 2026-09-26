CREATE OR REPLACE FUNCTION health_oauth_state_get(
    p_state TEXT
)
RETURNS TABLE (
    user_id BIGINT,
    expires_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN

    RETURN QUERY
    SELECT
        h.user_id,
        h.expires_at
    FROM health_oauth_states h
    WHERE h.state = p_state;

END;
$$;