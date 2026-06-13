package usda

import (
	"github.com/go-chi/chi/v5"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func RegisterRoutes(r chi.Router) {

	handler := NewHandler()

	r.Route("/usda", func(r chi.Router) {

		r.Get("/searchfoods", handler.SearchFoods)

	})
}
