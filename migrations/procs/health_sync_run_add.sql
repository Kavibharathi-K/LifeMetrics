CREATE OR REPLACE FUNCTION health_sync_run_add(
    p_user_id BIGINT,
    p_provider TEXT,
    p_sync_type TEXT
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    v_sync_run_id BIGINT;
BEGIN

    INSERT INTO health_sync_runs (
        user_id,
        provider,
        sync_type,
        started_at,
        status
    )
    VALUES (
        p_user_id,
        p_provider,
        p_sync_type,
        CURRENT_TIMESTAMP,
        'RUNNING'
    )
    RETURNING id INTO v_sync_run_id;

    RETURN v_sync_run_id;

END;
$$;