CREATE OR REPLACE FUNCTION list_foods()
RETURNS TABLE (
    food_id INT,
    name TEXT,
    measurement_type TEXT,
    base_quantity NUMERIC,
    calories NUMERIC,
    protein NUMERIC,
    carbs NUMERIC,
    fat NUMERIC,
    fiber NUMERIC,
    created_at TIMESTAMP
)
LANGUAGE sql
AS $$
    SELECT
        food_id,
        name,
        measurement_type,
        base_quantity,
        calories,
        protein,
        carbs,
        fat,
        fiber,
        created_at
    FROM foods
    ORDER BY name;
$$;
