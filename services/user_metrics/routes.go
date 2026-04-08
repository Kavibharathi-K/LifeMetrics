package user_metrics

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	r.Route("/usermetrics", func(r chi.Router) {

		r.Post("/addusermetrics", handler.CreateUserMetrics)

	})
}
func (h *Handler) CreateUserMetrics(w http.ResponseWriter, r *http.Request) {

	var req UserMetrics
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	calories := CalculateMaintenanceCalories(req.Age, req.Gender, req.HeightCm, req.WeightKg, req.ActivityLevel)
	req.MaintenanceCalories = calories

	proteinGoal, carbGoal, fatGoal := CalculateMacroGoals(calories, int(req.WeightKg))

	req.MaintenanceCalories = calories

	req.ProteinGoal = proteinGoal
	req.CarbGoal = carbGoal
	req.FatGoal = fatGoal
	_, err = h.Repo.CreateUserMetrics(r.Context(), req)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	json.NewEncoder(w).Encode(

		map[string]interface{}{
			"maintenance_calories": calories,
			"protein_goal": proteinGoal,
			"carb_goal":    carbGoal,
			"fat_goal":     fatGoal,
		},
	)
}
