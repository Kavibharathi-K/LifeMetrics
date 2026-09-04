package meal

import (
	"encoding/json"
	"net/http"
	"time"

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

	r.Route("/meals", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/gettodaynutrition", handler.GetTodayNutrition)
		r.Get("/getmealsbydate", handler.GetMealsByDate)
	})
}

func (h *Handler) GetTodayNutrition(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())

	if err != nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}
	var todayMeals []TodayNutrition
	todayMeals, err = h.Repo.GetTodayNutrition(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todayMeals)
}

func (h *Handler) GetMealsByDate(w http.ResponseWriter, r *http.Request) {

	userID, err := auth.GetUserID(r.Context())

	if err != nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	date := r.URL.Query().Get("date")

	if date == "" {

		date =
			time.Now().
				Format("2006-01-02")
	}

	response, err :=
		h.Repo.GetMealsByDate(
			r.Context(),
			userID,
			date,
		)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().
		Set(
			"Content-Type",
			"application/json",
		)

	json.NewEncoder(w).
		Encode(response)
}
