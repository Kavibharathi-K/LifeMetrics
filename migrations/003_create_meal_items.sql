CREATE TABLE meal_items (
    meal_item_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    meal_id INT NOT NULL,

    food_id INT NOT NULL,

    quantity DECIMAL(10,2) NOT NULL,

    calories DECIMAL(10,2) NOT NULL,
    protein DECIMAL(10,2) NOT NULL,
    carbs DECIMAL(10,2) NOT NULL,
    fat DECIMAL(10,2) NOT NULL,
    fiber DECIMAL(10,2) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_meal_items_meal
        FOREIGN KEY (meal_id)
        REFERENCES meals(meal_id)
        ON DELETE CASCADE
);