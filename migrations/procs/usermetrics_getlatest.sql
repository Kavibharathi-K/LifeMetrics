CREATE OR REPLACE FUNCTION usermetrics_getlatest()

RETURNS TABLE(

    id INT,

    age INT,
    gender TEXT,

    height_cm NUMERIC,
    weight_kg NUMERIC,

    activity_level TEXT,

    maintenance_calories INT,

    protein_goal INT,
    carb_goal INT,
    fat_goal INT,

    created_at TIMESTAMP

)

LANGUAGE sql

AS $$

    SELECT

        id,

        age,
        gender,

        height_cm,
        weight_kg,

        activity_level,

        maintenance_calories,

        protein_goal,
        carb_goal,
        fat_goal,

        created_at

    FROM user_metrics

    ORDER BY created_at DESC

    LIMIT 1;

$$;