CREATE TABLE foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    measurement_type TEXT NOT NULL DEFAULT 'grams',
    unit_name TEXT,

    calories_per_100g NUMERIC,
    protein_per_100g NUMERIC,
    carbs_per_100g NUMERIC,
    fat_per_100g NUMERIC,
    fiber_per_100g NUMERIC,

    calories_per_unit NUMERIC,
    protein_per_unit NUMERIC,
    carbs_per_unit NUMERIC,
    fat_per_unit NUMERIC,
    fiber_per_unit NUMERIC,

    created_at TIMESTAMP DEFAULT NOW()
);