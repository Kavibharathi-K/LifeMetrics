package usda

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) SearchFoods(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	foods, err := SearchFoods(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(foods)
}