CREATE TABLE workout_schedules (
    workout_schedule_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    day_of_week SMALLINT NOT NULL,
    workout_name VARCHAR(100) NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_workout_schedules_day_of_week
        CHECK (day_of_week BETWEEN 1 AND 7),

    CONSTRAINT chk_workout_schedules_name
        CHECK (LENGTH(TRIM(workout_name)) > 0),

    CONSTRAINT unique_workout_schedule_day
        UNIQUE (day_of_week)
);

CREATE TABLE workout_schedule_exercises (
    workout_schedule_exercise_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    workout_schedule_id INT NOT NULL,
    exercise_name VARCHAR(150) NOT NULL,
    exercise_order INT NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_workout_schedule_exercises_schedule
        FOREIGN KEY (workout_schedule_id)
        REFERENCES workout_schedules(workout_schedule_id)
        ON DELETE CASCADE,

    CONSTRAINT chk_workout_schedule_exercises_name
        CHECK (LENGTH(TRIM(exercise_name)) > 0),

    CONSTRAINT chk_workout_schedule_exercises_order
        CHECK (exercise_order > 0),

    CONSTRAINT unique_exercise_order_per_workout
        UNIQUE (workout_schedule_id, exercise_order)
);

CREATE INDEX idx_workout_schedule_exercises_schedule
    ON workout_schedule_exercises(workout_schedule_id);
