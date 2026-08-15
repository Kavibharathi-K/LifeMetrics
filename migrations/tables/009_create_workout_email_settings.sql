CREATE TABLE IF NOT EXISTS workout_email_settings (
    workout_email_settings_id SMALLINT PRIMARY KEY DEFAULT 1,
    email VARCHAR(254) NOT NULL DEFAULT '',
    email_time VARCHAR(5) NOT NULL DEFAULT '',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_workout_email_settings_singleton
        CHECK (workout_email_settings_id = 1),

    CONSTRAINT chk_workout_email_settings_values
        CHECK (
            (email = '' AND email_time = '')
            OR (
                LENGTH(TRIM(email)) > 0
                AND email_time ~ '^[0-2][0-9]:[0-5][0-9]$'
                AND email_time < '24:00'
            )
        )
);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'workout_schedules'
          AND column_name = 'email'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'workout_schedules'
          AND column_name = 'email_time'
    ) THEN
        INSERT INTO workout_email_settings (
            workout_email_settings_id,
            email,
            email_time
        )
        SELECT
            1,
            ws.email,
            ws.email_time
        FROM workout_schedules ws
        WHERE ws.email <> ''
          AND ws.email_time <> ''
        ORDER BY ws.updated_at DESC
        LIMIT 1
        ON CONFLICT (workout_email_settings_id) DO NOTHING;
    END IF;
END;
$$;

INSERT INTO workout_email_settings (
    workout_email_settings_id
)
SELECT 1
WHERE NOT EXISTS (
    SELECT 1
    FROM workout_email_settings
    WHERE workout_email_settings_id = 1
);
