package handlers

import (
	"auth-microservice/internal/crypto"
	jwtutil "auth-microservice/internal/jwt"
	"auth-microservice/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Mux  *chi.Mux
	Stor storage
	JWT  *jwtutil.Manager
}

type storage interface {
	UserByLogin(login string) (*models.User, error)
	AddUser(user *models.User) error
	UserExists() bool
}

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	user2, err := h.Stor.UserByLogin(user.Login)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if !crypto.CheckPwd([]byte(user.HashedPwd), []byte(user2.HashedPwd)) {
		fmt.Fprintln(w, "Incorrect")
	}

	token, err := h.JWT.Generate(int64(user.Id))
	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"access_token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User
	user.CreatedAt = time.Now()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error")
		log.Println(err)
	}
	err = h.Stor.AddUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Не удалось добавить пользователя"))
	}
	userJson, _ := json.Marshal(user)
	w.WriteHeader(http.StatusCreated)
	w.Write(userJson)
}

func AuthMiddleware(jwtM *jwtutil.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "no token", http.StatusUnauthorized)
				return
			}

			const prefix = "Bearer "
			if len(authHeader) < len(prefix) {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			tokenStr := authHeader[len(prefix):]

			claims, err := jwtM.Parse(tokenStr)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (h *Handlers) HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	json.NewEncoder(w).Encode(map[string]any{
		"user_id": userID,
	})
}

func (h *Handlers) SetupHandlers() {
	h.Mux.Post("/login", h.HandleLogin)
	h.Mux.Post("/register", h.HandleRegister)

	h.Mux.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(h.JWT))
		r.Get("/profile", h.HandleProfile)
	})
}
