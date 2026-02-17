package server

import (
	"auth-microservice/internal/database"
	"auth-microservice/internal/handlers"
	jwtutil "auth-microservice/internal/jwt"
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
	jwtManager := &jwtutil.Manager{
		Secret: []byte("super-secret-key"),
	}

	handler := handlers.Handlers{
		Mux:  chi.NewRouter(),
		Stor: db,
		JWT:  jwtManager,
	}
	return &Server{
		Handlers: handler,
	}
}
