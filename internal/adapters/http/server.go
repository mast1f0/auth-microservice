package http

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, r *chi.Mux) *Server {
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run() error {
	go func() {
		log.Println("Server running on", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http listen and serve error: %s\n", err)
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-shutdown:
		log.Printf("Received signal: %v\n", slog.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			if closeErr := s.httpServer.Close(); closeErr != nil {
				log.Printf("Force close error: %v\n", closeErr)
			}
			log.Println("graceful shutdown failed:", err)
			return err
		}

		log.Println("Server stopped gracefully")
		return nil
	}
}
