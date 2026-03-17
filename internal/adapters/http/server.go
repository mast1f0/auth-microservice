package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
}

func NewServer(r *chi.Mux) *Server {
	return &Server{
		router: r,
	}
}

func (s *Server) Run() error {
	port := 8081
	log.Printf("Listening on port %d", port)
	return http.ListenAndServe(":8081", s.router)
}
