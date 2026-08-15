CREATE OR REPLACE PROCEDURE add_meal_item(
    p_meal_type TEXT,
    p_meal_date DATE,
    p_food_name TEXT,
    p_quantity NUMERIC
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_meal_id INT;
    v_food_id INT;
    v_existing_item_id INT;
    v_base_quantity NUMERIC;
    v_calories NUMERIC;
    v_protein NUMERIC;
    v_carbs NUMERIC;
    v_fat NUMERIC;
    v_fiber NUMERIC;
    v_factor NUMERIC;
    v_item_calories NUMERIC;
    v_item_protein NUMERIC;
    v_item_carbs NUMERIC;
    v_item_fat NUMERIC;
    v_item_fiber NUMERIC;
BEGIN
    SELECT
        food_id,
        base_quantity,
        calories,
        protein,
        carbs,
        fat,
        fiber
    INTO
        v_food_id,
        v_base_quantity,
        v_calories,
        v_protein,
        v_carbs,
        v_fat,
        v_fiber
    FROM foods
    WHERE LOWER(name) = LOWER(p_food_name);

    IF v_food_id IS NULL THEN
        RAISE EXCEPTION 'food not found';
    END IF;

    v_factor := p_quantity / v_base_quantity;
    v_item_calories := v_calories * v_factor;
    v_item_protein := v_protein * v_factor;
    v_item_carbs := v_carbs * v_factor;
    v_item_fat := v_fat * v_factor;
    v_item_fiber := v_fiber * v_factor;

    SELECT meal_id
    INTO v_meal_id
    FROM meals
    WHERE meal_type = p_meal_type
      AND meal_date = p_meal_date;

    IF v_meal_id IS NULL THEN
        INSERT INTO meals (
            meal_type,
            meal_date
        )
        VALUES (
            p_meal_type,
            p_meal_date
        )
        RETURNING meal_id
        INTO v_meal_id;
    END IF;

    SELECT meal_item_id
    INTO v_existing_item_id
    FROM meal_items
    WHERE meal_id = v_meal_id
      AND food_id = v_food_id;

    IF v_existing_item_id IS NULL THEN
        INSERT INTO meal_items (
            meal_id,
            food_id,
            quantity,
            calories,
            protein,
            carbs,
            fat,
            fiber
        )
        VALUES (
            v_meal_id,
            v_food_id,
            p_quantity,
            v_item_calories,
            v_item_protein,
            v_item_carbs,
            v_item_fat,
            v_item_fiber
        );
    ELSE
        UPDATE meal_items
        SET
            quantity = quantity + p_quantity,
            calories = calories + v_item_calories,
            protein = protein + v_item_protein,
            carbs = carbs + v_item_carbs,
            fat = fat + v_item_fat,
            fiber = fiber + v_item_fiber
        WHERE meal_item_id = v_existing_item_id;
    END IF;

    UPDATE meals
    SET
        total_calories = total_calories + v_item_calories,
        total_protein = total_protein + v_item_protein,
        total_carbs = total_carbs + v_item_carbs,
        total_fat = total_fat + v_item_fat,
        total_fiber = total_fiber + v_item_fiber
    WHERE meal_id = v_meal_id;
END;
$$;
