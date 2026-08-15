CREATE OR REPLACE FUNCTION get_food_by_id(
    p_food_id INT
)
RETURNS TABLE (
    food_id INT,
    name VARCHAR(150),
    measurement_type VARCHAR(20),
    base_quantity DECIMAL(10,2),
    calories DECIMAL(10,2),
    protein DECIMAL(10,2),
    carbs DECIMAL(10,2),
    fat DECIMAL(10,2),
    fiber DECIMAL(10,2),
    created_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        f.food_id,
        f.name,
        f.measurement_type,
        f.base_quantity,
        f.calories,
        f.protein,
        f.carbs,
        f.fat,
        f.fiber,
        f.created_at
    FROM foods f
    WHERE f.food_id = p_food_id;
$$;
