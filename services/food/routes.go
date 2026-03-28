package food

import (
	"context"
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

	r.Route("/foods", func(r chi.Router) {

		r.Post("/", handler.CreateFood)
		r.Get("/", handler.ListFoods)

	})
}

func (h *Handler) CreateFood(w http.ResponseWriter, r *http.Request) {

	var food Food

	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.CreateFood(context.Background(), &food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(food)
}

func (h *Handler) ListFoods(w http.ResponseWriter, r *http.Request) {

	foods, err := h.Repo.ListFoods(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(foods)
}
