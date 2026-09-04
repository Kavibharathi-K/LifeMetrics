package user_metrics

import (
	"encoding/json"
	"net/http"

	"github.com/Kavibharathi-K/lifemetrics/services/auth"
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

		r.Use(auth.AuthMiddleware)

		r.Post("/addusermetrics", handler.CreateUserMetrics)
		r.Get("/getlatestusermetrics", handler.GetLatestUserMetrics)

	})
}
func (h *Handler) CreateUserMetrics(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())

	if err != nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return

	}
	var req UserMetrics
	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	req.UserId = userID

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
			"protein_goal":         proteinGoal,
			"carb_goal":            carbGoal,
			"fat_goal":             fatGoal,
		},
	)
}

func (h *Handler) GetLatestUserMetrics(

	w http.ResponseWriter,
	r *http.Request,

) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {

		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}
	data, err :=
		h.Repo.GetLatestUserMetrics(
			r.Context(),
			userID,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(w).
		Encode(data)
}
