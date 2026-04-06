create or replace procedure add_meal_item(
    p_meal_type text,
    p_meal_date date,
    p_food_name text,
    p_quantity numeric
)
language plpgsql
as $$

declare

    v_meal_id int;
    v_food_id int;
    v_base_quantity numeric;
    v_calories numeric;
    v_protein numeric;
    v_carbs numeric;
    v_fat numeric;
    v_fiber numeric;
    v_factor numeric;

begin

    /*
        find food by name
    */

    select
        food_id,
        base_quantity,
        calories,
        protein,
        carbs,
        fat,
        fiber
    into
        v_food_id,
        v_base_quantity,
        v_calories,
        v_protein,
        v_carbs,
        v_fat,
        v_fiber
    from foods
    where lower(name) = lower(p_food_name);

    if v_food_id is null then
        raise exception 'food not found';
    end if;

    /*
        calculate nutrition factor
    */

    v_factor := p_quantity / v_base_quantity;

    /*
        check if meal exists
    */

    select meal_id
    into v_meal_id
    from meals
    where meal_type = p_meal_type
    and meal_date = p_meal_date;

    /*
        create meal if not exists
    */

    if v_meal_id is null then
        insert into meals(
            meal_type,
            meal_date,
            total_calories,
            total_protein,
            total_carbs,
            total_fat,
            total_fiber
        )
        values (
            p_meal_type,
            p_meal_date,
            0,
            0,
            0,
            0,
            0
        )
        returning meal_id
        into v_meal_id;
    end if;

    /*
        insert meal item
    */

    insert into meal_items(
        meal_id,
        food_id,
        food_name,
        meal_type,
        quantity,
        calories,
        protein,
        carbs,
        fat,
        fiber
    )
    values(
        v_meal_id,
        v_food_id,
        p_food_name,
        p_meal_type,
        p_quantity,
        v_calories * v_factor,
        v_protein * v_factor,
        v_carbs * v_factor,
        v_fat * v_factor,
        v_fiber * v_factor
    );

    /*
        update meal totals
    */

    update meals
    set
        total_calories =
            total_calories + (v_calories * v_factor),

        total_protein =
            total_protein + (v_protein * v_factor),

        total_carbs =
            total_carbs + (v_carbs * v_factor),

        total_fat =
            total_fat + (v_fat * v_factor),

        total_fiber =
            total_fiber + (v_fiber * v_factor)
    where meal_id = v_meal_id;
end;
$$;