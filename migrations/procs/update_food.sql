CREATE OR REPLACE PROCEDURE update_food(
    p_food_id INT,
    p_name VARCHAR(150),
    p_measurement_type VARCHAR(20),
    p_base_quantity DECIMAL(10,2),
    p_calories DECIMAL(10,2),
    p_protein DECIMAL(10,2),
    p_carbs DECIMAL(10,2),
    p_fat DECIMAL(10,2),
    p_fiber DECIMAL(10,2)
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE foods
    SET
        name = p_name,
        measurement_type = p_measurement_type,
        base_quantity = p_base_quantity,
        calories = p_calories,
        protein = p_protein,
        carbs = p_carbs,
        fat = p_fat,
        fiber = p_fiber
    WHERE food_id = p_food_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'food with id % not found', p_food_id;
    END IF;
END;
$$;
