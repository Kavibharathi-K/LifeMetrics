CREATE OR REPLACE FUNCTION create_food(
    p_name TEXT,
    p_measurement_type TEXT,
    p_base_quantity NUMERIC,
    p_calories NUMERIC,
    p_protein NUMERIC,
    p_carbs NUMERIC,
    p_fat NUMERIC,
    p_fiber NUMERIC
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    v_food_id INT;
BEGIN
    INSERT INTO foods (
        name,
        measurement_type,
        base_quantity,
        calories,
        protein,
        carbs,
        fat,
        fiber
    )
    VALUES (
        p_name,
        p_measurement_type,
        p_base_quantity,
        p_calories,
        p_protein,
        p_carbs,
        p_fat,
        p_fiber
    )
    RETURNING food_id
    INTO v_food_id;

    RETURN v_food_id;
END;
$$;
