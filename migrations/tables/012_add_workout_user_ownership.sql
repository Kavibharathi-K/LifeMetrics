BEGIN;

-- ============================================================
-- Add user ownership to workout schedules
-- ============================================================

ALTER TABLE workout_schedules
ADD COLUMN user_id BIGINT;


-- ============================================================
-- Assign existing workout schedules to the current user
-- ============================================================

UPDATE workout_schedules
SET user_id = 4
WHERE user_id IS NULL;


-- ============================================================
-- user_id is required
-- ============================================================

ALTER TABLE workout_schedules
ALTER COLUMN user_id SET NOT NULL;


-- ============================================================
-- Foreign key to users
-- ============================================================

ALTER TABLE workout_schedules
ADD CONSTRAINT fk_workout_schedules_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;


-- ============================================================
-- Remove global day uniqueness
-- ============================================================

ALTER TABLE workout_schedules
DROP CONSTRAINT IF EXISTS unique_workout_schedule_day;


-- ============================================================
-- One workout per day per user
-- ============================================================

ALTER TABLE workout_schedules
ADD CONSTRAINT unique_user_workout_schedule_day
UNIQUE (user_id, day_of_week);


COMMIT;