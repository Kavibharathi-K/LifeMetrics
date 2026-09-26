CREATE OR REPLACE FUNCTION health_connection_get(
    p_user_id BIGINT,
    p_provider TEXT
)
RETURNS TABLE (
    provider_user_id TEXT,
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN

    RETURN QUERY
    SELECT
        h.provider_user_id,
        h.access_token,
        h.refresh_token,
        h.expires_at
    FROM health_connections h
    WHERE h.user_id = p_user_id
      AND h.provider = p_provider;

END;
$$;