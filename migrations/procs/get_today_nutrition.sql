CREATE OR REPLACE FUNCTION get_today_nutrition(
    p_user_id BIGINT,
    p_date DATE DEFAULT CURRENT_DATE
)
RETURNS TABLE (
    total_calories NUMERIC,
    total_protein NUMERIC,
    total_carbs NUMERIC,
    total_fat NUMERIC,
    total_fiber NUMERIC
)
LANGUAGE sql
AS $$
    SELECT
        COALESCE(SUM(m.total_calories), 0) AS total_calories,
        COALESCE(SUM(m.total_protein), 0) AS total_protein,
        COALESCE(SUM(m.total_carbs), 0) AS total_carbs,
        COALESCE(SUM(m.total_fat), 0) AS total_fat,
        COALESCE(SUM(m.total_fiber), 0) AS total_fiber
    FROM meals m
    WHERE m.user_id = p_user_id
      AND m.meal_date = p_date;
$$;