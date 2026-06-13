package workout

import "time"

type WorkoutSchedule struct {
	WorkoutScheduleId int       `db:"workout_schedule_id" json:"workout_schedule_id"`
	DayOfWeek         int       `db:"day_of_week" json:"day_of_week"`
	WorkoutName       string    `db:"workout_name" json:"workout_name"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type WorkoutScheduleExercises struct {
	WorkoutScheduleExerciseId int       `db:"workout_schedule_exercise_id" json:"workout_schedule_exercise_id"`
	WorkoutScheduleId         int       `db:"workout_schedule_id" json:"workout_schedule_id"`
	ExerciseName              string    `db:"exercise_name" json:"exercise_name"`
	ExerciseOrder             int       `db:"exercise_order" json:"exercise_order"`
	CreatedAt                 time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at" json:"updated_at"`
}

type CreateWorkoutScheduleRequest struct {
	DayOfWeek   int      `json:"day_of_week"`
	WorkoutName string   `json:"workout_name"`
	Exercises   []string `json:"exercises"`
}

type CreateWorkoutScheduleResponse struct {
	WorkoutScheduleId int      `json:"workout_schedule_id"`
	DayOfWeek         int      `json:"day_of_week"`
	WorkoutName       string   `json:"workout_name"`
	Exercises         []string `json:"exercises"`
}

type UpdateWorkoutScheduleRequest struct {
	DayOfWeek   int      `json:"day_of_week"`
	WorkoutName string   `json:"workout_name"`
	Exercises   []string `json:"exercises"`
}

type UpdateWorkoutScheduleResponse struct {
	WorkoutScheduleId int      `json:"workout_schedule_id"`
	DayOfWeek         int      `json:"day_of_week"`
	WorkoutName       string   `json:"workout_name"`
	Exercises         []string `json:"exercises"`
}

type CreateWorkoutScheduleExerciseRequest struct {
	ExerciseName  string `json:"exercise_name"`
	ExerciseOrder int    `json:"exercise_order,omitempty"`
}

type CreateWorkoutScheduleExerciseResponse struct {
	WorkoutScheduleExerciseId int    `json:"workout_schedule_exercise_id"`
	WorkoutScheduleId         int    `json:"workout_schedule_id"`
	ExerciseName              string `json:"exercise_name"`
	ExerciseOrder             int    `json:"exercise_order"`
}

type WorkoutScheduleExerciseResponse struct {
	WorkoutScheduleExerciseId int    `json:"workout_schedule_exercise_id"`
	ExerciseName              string `json:"exercise_name"`
	ExerciseOrder             int    `json:"exercise_order"`
}

type WorkoutScheduleResponse struct {
	WorkoutScheduleId int                               `json:"workout_schedule_id"`
	DayOfWeek         int                               `json:"day_of_week"`
	WorkoutName       string                            `json:"workout_name"`
	Exercises         []WorkoutScheduleExerciseResponse `json:"exercises"`
	CreatedAt         time.Time                         `json:"created_at"`
	UpdatedAt         time.Time                         `json:"updated_at"`
}
