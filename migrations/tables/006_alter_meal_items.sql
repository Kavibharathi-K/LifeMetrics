alter table meal_items

add constraint unique_meal_food

unique(meal_id, food_id);