package server

import (
	"auth-microservice/internal/database"
	"auth-microservice/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Handlers handlers.Handlers
}

func (srv *Server) RunServer() {
	srv.Handlers.SetupHandlers()
	mux := srv.Handlers.Mux
	http.ListenAndServe(":8080", mux)
}

func NewServer(db *database.Database) *Server {
	handler := handlers.Handlers{
		Mux:  chi.NewRouter(),
		Stor: db,
	}
	return &Server{
		Handlers: handler,
	}
}
