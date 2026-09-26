CREATE OR REPLACE PROCEDURE health_connection_upsert(
    p_user_id BIGINT,
    p_provider TEXT,
    p_provider_user_id TEXT,
    p_access_token TEXT,
    p_refresh_token TEXT,
    p_expires_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN

    INSERT INTO health_connections (
        user_id,
        provider,
        provider_user_id,
        access_token,
        refresh_token,
        expires_at
    )
    VALUES (
        p_user_id,
        p_provider,
        p_provider_user_id,
        p_access_token,
        p_refresh_token,
        p_expires_at
    )
    ON CONFLICT (user_id, provider)
    DO UPDATE SET
        provider_user_id = EXCLUDED.provider_user_id,
        access_token = EXCLUDED.access_token,
        refresh_token = CASE
            WHEN EXCLUDED.refresh_token IS NOT NULL
                 AND EXCLUDED.refresh_token <> ''
            THEN EXCLUDED.refresh_token
            ELSE health_connections.refresh_token
        END,
        expires_at = EXCLUDED.expires_at,
        updated_at = CURRENT_TIMESTAMP;

END;
$$;