CREATE OR REPLACE FUNCTION get_meals_by_date(p_date DATE)

RETURNS TABLE (

    meal_id INT,
    meal_type VARCHAR,

    total_calories NUMERIC,
    total_protein NUMERIC,
    total_carbs NUMERIC,
    total_fat NUMERIC,
    total_fiber NUMERIC,

    food_id INT,
    food_name VARCHAR,
    measurement_type VARCHAR,

    quantity NUMERIC,
    calories NUMERIC,
    protein NUMERIC,
    carbs NUMERIC,
    fat NUMERIC,
    fiber NUMERIC

)

LANGUAGE SQL

AS $$

SELECT

    m.meal_id,
    m.meal_type,

    m.total_calories,
    m.total_protein,
    m.total_carbs,
    m.total_fat,
    m.total_fiber,

    f.food_id,
    f.name,
    f.measurement_type,

    mi.quantity,
    mi.calories,
    mi.protein,
    mi.carbs,
    mi.fat,
    mi.fiber

FROM meals m

LEFT JOIN meal_items mi
    ON mi.meal_id = m.meal_id

LEFT JOIN foods f
    ON f.food_id = mi.food_id

WHERE m.meal_date = p_date

ORDER BY m.meal_type;

$$;