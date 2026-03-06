package calories

import (
	"context"
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{Repo: repo}
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
