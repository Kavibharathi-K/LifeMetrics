CREATE TABLE meals (
    meal_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    meal_date DATE NOT NULL,

    meal_type VARCHAR(20) NOT NULL,
    -- breakfast, lunch, dinner, snack

    total_calories DECIMAL(10,2) DEFAULT 0,
    total_protein DECIMAL(10,2) DEFAULT 0,
    total_carbs DECIMAL(10,2) DEFAULT 0,
    total_fat DECIMAL(10,2) DEFAULT 0,
    total_fiber DECIMAL(10,2) DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);