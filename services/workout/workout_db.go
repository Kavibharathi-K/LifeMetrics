package workout

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateWorkoutSchedule(ctx context.Context, schedule *WorkoutSchedule, exercises []string) error {
	query := `
		SELECT create_workout_schedule($1, $2, $3)
	`

	return r.DB.QueryRow(
		ctx,
		query,
		schedule.DayOfWeek,
		schedule.WorkoutName,
		exercises,
	).Scan(&schedule.WorkoutScheduleId)
}

func (r *Repository) CreateWorkoutScheduleExercise(ctx context.Context, exercise *WorkoutScheduleExercises) error {
	query := `
		SELECT *
		FROM create_workout_schedule_exercise($1, $2, $3)
	`

	return r.DB.QueryRow(
		ctx,
		query,
		exercise.WorkoutScheduleId,
		exercise.ExerciseName,
		exercise.ExerciseOrder,
	).Scan(
		&exercise.WorkoutScheduleExerciseId,
		&exercise.ExerciseOrder,
	)
}

func (r *Repository) ListWorkoutSchedules(ctx context.Context) ([]WorkoutScheduleResponse, error) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM list_workout_schedules()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]WorkoutScheduleResponse, 0)

	for rows.Next() {
		var schedule WorkoutScheduleResponse
		var exercisesJSON []byte

		if err := rows.Scan(
			&schedule.WorkoutScheduleId,
			&schedule.DayOfWeek,
			&schedule.WorkoutName,
			&exercisesJSON,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(exercisesJSON, &schedule.Exercises); err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *Repository) GetWorkoutScheduleByDay(ctx context.Context, dayOfWeek int) (WorkoutScheduleResponse, error) {
	var schedule WorkoutScheduleResponse
	var exercisesJSON []byte

	err := r.DB.QueryRow(
		ctx,
		`SELECT * FROM get_workout_schedule_by_day($1)`,
		dayOfWeek,
	).Scan(
		&schedule.WorkoutScheduleId,
		&schedule.DayOfWeek,
		&schedule.WorkoutName,
		&exercisesJSON,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	if err != nil {
		return WorkoutScheduleResponse{}, err
	}

	if err := json.Unmarshal(exercisesJSON, &schedule.Exercises); err != nil {
		return WorkoutScheduleResponse{}, err
	}

	return schedule, nil
}

func (r *Repository) UpdateWorkoutSchedule(
	ctx context.Context,
	schedule *WorkoutSchedule,
	exercises []string,
) error {
	_, err := r.DB.Exec(
		ctx,
		`CALL update_workout_schedule($1, $2, $3, $4)`,
		schedule.WorkoutScheduleId,
		schedule.DayOfWeek,
		schedule.WorkoutName,
		exercises,
	)

	return err
}

func (r *Repository) DeleteWorkoutSchedule(
	ctx context.Context,
	workoutScheduleID int,
) (bool, error) {
	var deleted bool

	err := r.DB.QueryRow(
		ctx,
		`SELECT delete_workout_schedule($1)`,
		workoutScheduleID,
	).Scan(&deleted)

	return deleted, err
}

func (r *Repository) GetWorkoutEmailSettings(ctx context.Context) (WorkoutEmailSettings, error) {
	var settings WorkoutEmailSettings

	err := r.DB.QueryRow(
		ctx,
		`
			SELECT email, email_time, updated_at
			FROM workout_email_settings
			WHERE workout_email_settings_id = 1
		`,
	).Scan(
		&settings.Email,
		&settings.EmailTime,
		&settings.UpdatedAt,
	)

	return settings, err
}

func (r *Repository) UpdateWorkoutEmailSettings(
	ctx context.Context,
	settings *WorkoutEmailSettings,
) error {
	return r.DB.QueryRow(
		ctx,
		`
			INSERT INTO workout_email_settings (
				workout_email_settings_id,
				email,
				email_time
			)
			VALUES (1, $1, $2)
			ON CONFLICT (workout_email_settings_id)
			DO UPDATE SET
				email = EXCLUDED.email,
				email_time = EXCLUDED.email_time,
				updated_at = CURRENT_TIMESTAMP
			RETURNING updated_at
		`,
		settings.Email,
		settings.EmailTime,
	).Scan(&settings.UpdatedAt)
}
