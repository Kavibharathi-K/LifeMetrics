package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r chi.Router, db *pgxpool.Pool) {

	repo := NewRepository(db)
	handler := NewHandler(repo)

	r.Route("/auth", func(r chi.Router) {

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}