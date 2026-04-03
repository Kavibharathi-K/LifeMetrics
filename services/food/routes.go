package food

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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

		r.Post("/addfood", handler.CreateFood)
		r.Get("/listfoods", handler.ListFoods)
		r.Get("/getfood/{food_id}", handler.GetFoodByID)
		r.Post("/updatefood/{food_id}", handler.UpdateFood)

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

func (h *Handler) GetFoodByID(w http.ResponseWriter, r *http.Request) {

	foodIDStr := chi.URLParam(r, "food_id")

	foodID, err := strconv.Atoi(foodIDStr)
	if err != nil {
		http.Error(w, "invalid food_id", http.StatusBadRequest)
		return
	}

	food, err := h.Repo.GetFoodByID(context.Background(), foodID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(food)
}

func (h *Handler) UpdateFood(w http.ResponseWriter, r *http.Request) {

	foodIDstr := chi.URLParam(r, "food_id")

	foodID, err := strconv.Atoi(foodIDstr)
	if err != nil {
		http.Error(w, "invalid food_id", http.StatusBadRequest)
		return
	}

	var food Food
	food.FoodID = foodID
	err = json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Repo.UpdateFood(context.Background(), foodID, &food)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(food)
}
