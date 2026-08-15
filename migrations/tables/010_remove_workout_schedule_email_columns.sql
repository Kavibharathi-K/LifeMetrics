DROP FUNCTION IF EXISTS list_workout_schedules();
DROP FUNCTION IF EXISTS get_workout_schedule_by_day(SMALLINT);
DROP FUNCTION IF EXISTS create_workout_schedule(SMALLINT, TEXT, TEXT, TEXT, TEXT[]);
DROP PROCEDURE IF EXISTS update_workout_schedule(INT, SMALLINT, TEXT, TEXT, TEXT, TEXT[]);

ALTER TABLE workout_schedules
    DROP CONSTRAINT IF EXISTS chk_workout_schedules_email_settings,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS email_time;
