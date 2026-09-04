package meal_item

import (
	"encoding/json"
	"log"
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
	r.Route("/mealitems", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/addmealitem", handler.AddMealItem)
		r.Get("/getMealItem", handler.GetMealItem)
	})

}

func (h *Handler) AddMealItem(w http.ResponseWriter, r *http.Request) {

	var meal MealItem

	err := json.NewDecoder(r.Body).Decode(&meal)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	userID, err := auth.GetUserID(r.Context())

	if err != nil {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	log.Printf(
		"MealItem received: %+v\n",
		meal,
	)

	today := time.Now().Format("2006-01-02")

	err = h.Repo.AddMealItem(
		r.Context(),
		userID,
		meal.MealType,
		today,
		meal.FoodName,
		meal.Quantity,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"message":   "Meal item added successfully",
			"meal_type": meal.MealType,
			"food_name": meal.FoodName,
			"quantity":  meal.Quantity,
		},
	)
}

func (h *Handler) GetMealItem(
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

	items, err := h.Repo.GetMealItems(
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

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(items)
}
