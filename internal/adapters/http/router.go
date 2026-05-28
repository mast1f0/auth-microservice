package http

import (
	"auth-microservice/internal/adapters/http/handlers"
	"auth-microservice/internal/adapters/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(h *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
	}))
	r.Post("/login", h.HandleLogin)
	r.Post("/register", h.HandleRegister)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(h.JWT))
		r.Get("/profile", h.HandleProfile)
		r.Delete("/delete", h.DeleteUser)
		r.Get("/users", h.GetAllUsers)
		r.Put("/users", h.UpdateUser)
	})
	return r
}
