create or replace procedure update_food(
    p_food_id int,
    p_name varchar(150),
    p_measurement_type varchar(20),
    p_base_quantity decimal(10,2),
    p_calories decimal(10,2),
    p_protein decimal(10,2),
    p_carbs decimal(10,2),
    p_fat decimal(10,2),
    p_fiber decimal(10,2)
)

language plpgsql
as $$ 
begin 
    update foods
    set
        name = p_name,
        measurement_type = p_measurement_type,
        base_quantity = p_base_quantity,
        calories = p_calories,
        protein = p_protein,
        carbs = p_carbs,
        fat = p_fat,
        fiber = p_fiber
    where food_id = p_food_id;

    if not found then 
        raise exception 'food with id % not found', p_food_id;
    end if;

end;
$$;