package health

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router) {
	r.HandleFunc("/health/google/connect", GoogleConnectHandler)
	r.HandleFunc("/health/google/callback", GoogleCallbackHandler)
}
