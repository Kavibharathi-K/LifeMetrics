CREATE OR REPLACE PROCEDURE health_oauth_state_delete(
    p_state TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN

    DELETE FROM health_oauth_states
    WHERE state = p_state;

END;
$$;