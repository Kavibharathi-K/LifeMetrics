CREATE TYPE meal_type AS ENUM ('breakfast', 'lunch', 'dinner', 'snack');

CREATE TABLE meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_date DATE NOT NULL,
    meal_type meal_type NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_meals_user_date
ON meals(user_id, meal_date);