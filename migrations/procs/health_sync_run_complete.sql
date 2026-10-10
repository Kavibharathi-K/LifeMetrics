CREATE OR REPLACE PROCEDURE health_sync_run_complete(
    p_sync_run_id BIGINT,
    p_status TEXT,
    p_error_message TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN

    UPDATE health_sync_runs
    SET
        completed_at = CURRENT_TIMESTAMP,
        status = p_status,
        error_message = p_error_message
    WHERE id = p_sync_run_id;

END;
$$;