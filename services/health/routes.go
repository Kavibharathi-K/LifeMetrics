package health

import (
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

	r.Route("/health/google", func(r chi.Router) {

		r.With(auth.AuthMiddleware).Get(
			"/connect",
			handler.GoogleConnectHandler,
		)

		r.Get(
			"/callback",
			handler.GoogleCallbackHandler,
		)

		r.With(auth.AuthMiddleware).Get(
			"/test-token",
			handler.TestTokenHandler,
		)

		r.With(auth.AuthMiddleware).Post("/sync-daily-activity", handler.SyncDailyActivity)

		r.With(auth.AuthMiddleware).Post("/sync-daily-vitals", handler.SyncDailyVitalsHandler)
	})
}
