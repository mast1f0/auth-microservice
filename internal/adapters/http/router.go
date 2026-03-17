package http

import (
	"auth-microservice/internal/adapters/http/handlers"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/login", h.HandleLogin)
	r.Post("/register", h.HandleRegister)

	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(h.JWT))
		r.Get("/profile", h.HandleProfile)
	})
	return r
}
