CREATE OR REPLACE FUNCTION usermetrics_getlatest(
    p_user_id BIGINT
)
RETURNS TABLE (
    id INT,
    user_id BIGINT,
    age INT,
    gender TEXT,
    height_cm NUMERIC,
    weight_kg NUMERIC,
    activity_level TEXT,
    maintenance_calories INT,
    protein_goal INT,
    carb_goal INT,
    fat_goal INT,
    created_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        id,
        user_id,
        age,
        gender,
        height_cm,
        weight_kg,
        activity_level,
        maintenance_calories,
        protein_goal,
        carb_goal,
        fat_goal,
        created_at
    FROM user_metrics
    WHERE user_id = p_user_id
    LIMIT 1;
$$;