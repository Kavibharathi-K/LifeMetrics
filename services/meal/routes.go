package meal

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

	r.Route("/meals", func(r chi.Router) {
		r.Get("/gettodaynutrition", handler.GetTodayNutrition)
	})
}

func (h *Handler) GetTodayNutrition(w http.ResponseWriter, r *http.Request) {
	var todayMeals []TodayNutrition
	todayMeals, err := h.Repo.GetTodayNutrition(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todayMeals)
}
