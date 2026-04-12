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

    v_existing_item_id int;

    v_base_quantity numeric;

    v_calories numeric;

    v_protein numeric;

    v_carbs numeric;

    v_fat numeric;

    v_fiber numeric;

    v_factor numeric;

    v_item_calories numeric;

    v_item_protein numeric;

    v_item_carbs numeric;

    v_item_fat numeric;

    v_item_fiber numeric;

begin

    /*
        find food
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
        calculate nutrition
    */

    v_factor := p_quantity / v_base_quantity;

    v_item_calories := v_calories * v_factor;

    v_item_protein := v_protein * v_factor;

    v_item_carbs := v_carbs * v_factor;

    v_item_fat := v_fat * v_factor;

    v_item_fiber := v_fiber * v_factor;


    /*
        get meal
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
        values(

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
        check if food already exists in meal
    */

    select meal_item_id
    into v_existing_item_id
    from meal_items
    where meal_id = v_meal_id
    and food_id = v_food_id;


    /*
        update if exists
    */

    if v_existing_item_id is not null then

        update meal_items
        set

            quantity = quantity + p_quantity,

            calories = calories + v_item_calories,

            protein = protein + v_item_protein,

            carbs = carbs + v_item_carbs,

            fat = fat + v_item_fat,

            fiber = fiber + v_item_fiber

        where meal_item_id = v_existing_item_id;

    else

        /*
            insert new food row
        */

        insert into meal_items(

            meal_id,
            food_id,
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
            p_quantity,
            v_item_calories,
            v_item_protein,
            v_item_carbs,
            v_item_fat,
            v_item_fiber

        );

    end if;


    /*
        update meal totals
    */

    update meals
    set

        total_calories =
            total_calories + v_item_calories,

        total_protein =
            total_protein + v_item_protein,

        total_carbs =
            total_carbs + v_item_carbs,

        total_fat =
            total_fat + v_item_fat,

        total_fiber =
            total_fiber + v_item_fiber

    where meal_id = v_meal_id;


end;
$$;