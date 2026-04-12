CREATE OR REPLACE FUNCTION get_today_nutrition(p_date date DEFAULT CURRENT_DATE)
RETURNS TABLE (
    total_calories numeric,
    total_protein  numeric,
    total_carbs    numeric,
    total_fat      numeric,
    total_fiber    numeric
)
AS
$$
BEGIN
    RETURN QUERY
    SELECT
        COALESCE(SUM(m.total_calories), 0) AS total_calories,
        COALESCE(SUM(m.total_protein), 0)  AS total_protein,
        COALESCE(SUM(m.total_carbs), 0)    AS total_carbs,
        COALESCE(SUM(m.total_fat), 0)      AS total_fat,
        COALESCE(SUM(m.total_fiber), 0)    AS total_fiber
    FROM meals m
    WHERE m.meal_date = p_date;
END;
$$
LANGUAGE plpgsql;