package workout

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
}

func RegisterRoutes(r chi.Router, db *pgxpool.Pool) {
	repo := NewRepository(db)
	handler := NewHandler(repo)

	r.Route("/workouts", func(r chi.Router) {
		r.Get("/schedules", handler.ListWorkoutSchedules)
		r.Get("/schedules/{day_of_week}", handler.GetWorkoutScheduleByDay)
		r.Post("/schedules", handler.CreateWorkoutSchedule)
		r.Put("/schedules/{workout_schedule_id}", handler.UpdateWorkoutSchedule)
		r.Delete("/schedules/{workout_schedule_id}", handler.DeleteWorkoutSchedule)
		r.Post("/schedules/{workout_schedule_id}/exercises", handler.CreateWorkoutScheduleExercise)
	})
}

// List whole week workout schedule.
func (h *Handler) ListWorkoutSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.Repo.ListWorkoutSchedules(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list workouts")
		return
	}

	writeJSON(w, http.StatusOK, schedules)
}

// List the workout for the day.
func (h *Handler) GetWorkoutScheduleByDay(w http.ResponseWriter, r *http.Request) {
	dayOfWeek, err := strconv.Atoi(chi.URLParam(r, "day_of_week"))
	if err != nil || dayOfWeek < 1 || dayOfWeek > 7 {
		writeError(w, http.StatusBadRequest, "day_of_week must be between 1 and 7")
		return
	}

	schedule, err := h.Repo.GetWorkoutScheduleByDay(
		r.Context(),
		dayOfWeek,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "workout schedule not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get workout")
		return
	}

	writeJSON(w, http.StatusOK, schedule)
}

// Creates an entire new workout schedule.
func (h *Handler) CreateWorkoutSchedule(w http.ResponseWriter, r *http.Request) {
	var request CreateWorkoutScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	request.WorkoutName = strings.TrimSpace(request.WorkoutName)

	if request.DayOfWeek < 1 || request.DayOfWeek > 7 {
		writeError(w, http.StatusBadRequest, "day_of_week must be between 1 and 7")
		return
	}

	if request.WorkoutName == "" {
		writeError(w, http.StatusBadRequest, "workout_name is required")
		return
	}

	for index, exercise := range request.Exercises {
		request.Exercises[index] = strings.TrimSpace(exercise)
		if request.Exercises[index] == "" {
			writeError(w, http.StatusBadRequest, "exercise names cannot be empty")
			return
		}
	}

	schedule := WorkoutSchedule{
		DayOfWeek:   request.DayOfWeek,
		WorkoutName: request.WorkoutName,
	}

	if err := h.Repo.CreateWorkoutSchedule(
		r.Context(),
		&schedule,
		request.Exercises,
	); err != nil {
		writeDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateWorkoutScheduleResponse{
		WorkoutScheduleId: schedule.WorkoutScheduleId,
		DayOfWeek:         schedule.DayOfWeek,
		WorkoutName:       schedule.WorkoutName,
		Exercises:         request.Exercises,
	})
}

// Update workout schedule for the day
func (h *Handler) UpdateWorkoutSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := strconv.Atoi(
		chi.URLParam(r, "workout_schedule_id"),
	)
	if err != nil || scheduleID < 1 {
		writeError(w, http.StatusBadRequest, "invalid workout_schedule_id")
		return
	}

	var request UpdateWorkoutScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	request.WorkoutName = strings.TrimSpace(request.WorkoutName)

	if request.DayOfWeek < 1 || request.DayOfWeek > 7 {
		writeError(w, http.StatusBadRequest, "day_of_week must be between 1 and 7")
		return
	}

	if request.WorkoutName == "" {
		writeError(w, http.StatusBadRequest, "workout_name is required")
		return
	}

	for index, exercise := range request.Exercises {
		request.Exercises[index] = strings.TrimSpace(exercise)
		if request.Exercises[index] == "" {
			writeError(w, http.StatusBadRequest, "exercise names cannot be empty")
			return
		}
	}

	schedule := WorkoutSchedule{
		WorkoutScheduleId: scheduleID,
		DayOfWeek:         request.DayOfWeek,
		WorkoutName:       request.WorkoutName,
	}

	if err := h.Repo.UpdateWorkoutSchedule(
		r.Context(),
		&schedule,
		request.Exercises,
	); err != nil {
		writeDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, UpdateWorkoutScheduleResponse{
		WorkoutScheduleId: schedule.WorkoutScheduleId,
		DayOfWeek:         schedule.DayOfWeek,
		WorkoutName:       schedule.WorkoutName,
		Exercises:         request.Exercises,
	})
}

// Deletes the workout schedule
func (h *Handler) DeleteWorkoutSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := strconv.Atoi(
		chi.URLParam(r, "workout_schedule_id"),
	)
	if err != nil || scheduleID < 1 {
		writeError(w, http.StatusBadRequest, "invalid workout_schedule_id")
		return
	}

	deleted, err := h.Repo.DeleteWorkoutSchedule(r.Context(), scheduleID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete workout")
		return
	}

	if !deleted {
		writeError(w, http.StatusNotFound, "workout schedule not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Add one exercise to a workout that already exists.
func (h *Handler) CreateWorkoutScheduleExercise(w http.ResponseWriter, r *http.Request) {
	scheduleID, err := strconv.Atoi(
		chi.URLParam(r, "workout_schedule_id"),
	)
	if err != nil || scheduleID < 1 {
		writeError(w, http.StatusBadRequest, "invalid workout_schedule_id")
		return
	}

	var request CreateWorkoutScheduleExerciseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request")
		return
	}

	request.ExerciseName = strings.TrimSpace(request.ExerciseName)

	if request.ExerciseName == "" {
		writeError(w, http.StatusBadRequest, "exercise_name is required")
		return
	}

	if request.ExerciseOrder < 0 {
		writeError(w, http.StatusBadRequest, "exercise_order cannot be negative")
		return
	}

	exercise := WorkoutScheduleExercises{
		WorkoutScheduleId: scheduleID,
		ExerciseName:      request.ExerciseName,
		ExerciseOrder:     request.ExerciseOrder,
	}

	if err := h.Repo.CreateWorkoutScheduleExercise(
		r.Context(),
		&exercise,
	); err != nil {
		writeDatabaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateWorkoutScheduleExerciseResponse{
		WorkoutScheduleExerciseId: exercise.WorkoutScheduleExerciseId,
		WorkoutScheduleId:         exercise.WorkoutScheduleId,
		ExerciseName:              exercise.ExerciseName,
		ExerciseOrder:             exercise.ExerciseOrder,
	})
}

func writeDatabaseError(w http.ResponseWriter, err error) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			writeError(
				w,
				http.StatusConflict,
				"a workout or exercise already uses that day or order",
			)
			return
		case "P0002":
			writeError(w, http.StatusNotFound, "workout schedule not found")
			return
		case "23514":
			writeError(w, http.StatusBadRequest, "workout data is invalid")
			return
		}
	}

	writeError(w, http.StatusInternalServerError, "workout database operation failed")
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
