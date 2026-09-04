CREATE OR REPLACE FUNCTION user_metrics_add(
    p_user_id BIGINT,
    p_age INT,
    p_gender TEXT,
    p_height_cm NUMERIC,
    p_weight_kg NUMERIC,
    p_activity_level TEXT,
    p_maintenance_calories INT,
    p_protein_goal INT,
    p_carb_goal INT,
    p_fat_goal INT
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    v_id INT;
BEGIN
    INSERT INTO user_metrics (
        user_id,
        age,
        gender,
        height_cm,
        weight_kg,
        activity_level,
        maintenance_calories,
        protein_goal,
        carb_goal,
        fat_goal
    )
    VALUES (
        p_user_id,
        p_age,
        p_gender,
        p_height_cm,
        p_weight_kg,
        p_activity_level,
        p_maintenance_calories,
        p_protein_goal,
        p_carb_goal,
        p_fat_goal
    )
    ON CONFLICT (user_id)
    DO UPDATE SET
        age = EXCLUDED.age,
        gender = EXCLUDED.gender,
        height_cm = EXCLUDED.height_cm,
        weight_kg = EXCLUDED.weight_kg,
        activity_level = EXCLUDED.activity_level,
        maintenance_calories = EXCLUDED.maintenance_calories,
        protein_goal = EXCLUDED.protein_goal,
        carb_goal = EXCLUDED.carb_goal,
        fat_goal = EXCLUDED.fat_goal
    RETURNING id
    INTO v_id;

    RETURN v_id;
END;
$$;