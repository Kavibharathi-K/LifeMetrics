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

func (r *Repository) CreateWorkoutSchedule(
	ctx context.Context,
	userID int,
	schedule *WorkoutSchedule,
	exercises []string,
) error {

	query := `
		SELECT create_workout_schedule($1, $2, $3, $4)
	`

	return r.DB.QueryRow(
		ctx,
		query,
		userID,
		schedule.DayOfWeek,
		schedule.WorkoutName,
		exercises,
	).Scan(&schedule.WorkoutScheduleId)
}

func (r *Repository) CreateWorkoutScheduleExercise(
	ctx context.Context,
	userID int,
	exercise *WorkoutScheduleExercises,
) error {

	query := `
		SELECT *
		FROM create_workout_schedule_exercise($1, $2, $3, $4)
	`

	return r.DB.QueryRow(
		ctx,
		query,
		userID,
		exercise.WorkoutScheduleId,
		exercise.ExerciseName,
		exercise.ExerciseOrder,
	).Scan(
		&exercise.WorkoutScheduleExerciseId,
		&exercise.ExerciseOrder,
	)
}

func (r *Repository) ListWorkoutSchedules(
	ctx context.Context,
	userID int,
) ([]WorkoutScheduleResponse, error) {

	rows, err := r.DB.Query(
		ctx,
		`SELECT * FROM list_workout_schedules($1)`,
		userID,
	)

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

		if err := json.Unmarshal(
			exercisesJSON,
			&schedule.Exercises,
		); err != nil {
			return nil, err
		}

		schedules = append(
			schedules,
			schedule,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *Repository) GetWorkoutScheduleByDay(
	ctx context.Context,
	userID int,
	dayOfWeek int,
) (WorkoutScheduleResponse, error) {

	var schedule WorkoutScheduleResponse
	var exercisesJSON []byte

	err := r.DB.QueryRow(
		ctx,
		`SELECT * FROM get_workout_schedule_by_day($1, $2)`,
		userID,
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

	if err := json.Unmarshal(
		exercisesJSON,
		&schedule.Exercises,
	); err != nil {
		return WorkoutScheduleResponse{}, err
	}

	return schedule, nil
}

func (r *Repository) UpdateWorkoutSchedule(
	ctx context.Context,
	userID int,
	schedule *WorkoutSchedule,
	exercises []string,
) error {

	_, err := r.DB.Exec(
		ctx,
		`CALL update_workout_schedule($1, $2, $3, $4, $5)`,
		schedule.WorkoutScheduleId,
		userID,
		schedule.DayOfWeek,
		schedule.WorkoutName,
		exercises,
	)

	return err
}

func (r *Repository) DeleteWorkoutSchedule(
	ctx context.Context,
	userID int,
	workoutScheduleID int,
) (bool, error) {

	var deleted bool

	err := r.DB.QueryRow(
		ctx,
		`SELECT delete_workout_schedule($1, $2)`,
		workoutScheduleID,
		userID,
	).Scan(&deleted)

	return deleted, err
}

func (r *Repository) GetWorkoutEmailUsers(
	ctx context.Context,
) ([]int, error) {

	rows, err := r.DB.Query(
		ctx,
		`SELECT * FROM get_workout_email_users()`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userIDs := make([]int, 0)

	for rows.Next() {

		var userID int

		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}

		userIDs = append(
			userIDs,
			userID,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}

func (r *Repository) GetWorkoutEmailSettings(
	ctx context.Context,
	userID int,
) (WorkoutEmailSettings, error) {

	var settings WorkoutEmailSettings

	err := r.DB.QueryRow(
		ctx,
		`SELECT * FROM get_workout_email_settings($1)`,
		userID,
	).Scan(
		&settings.Email,
		&settings.EmailTime,
		&settings.UpdatedAt,
	)

	return settings, err
}

func (r *Repository) UpdateWorkoutEmailSettings(
	ctx context.Context,
	userID int,
	settings *WorkoutEmailSettings,
) error {

	return r.DB.QueryRow(
		ctx,
		`SELECT update_workout_email_settings($1, $2, $3)`,
		userID,
		settings.Email,
		settings.EmailTime,
	).Scan(&settings.UpdatedAt)
}