package handlers

import (
	"auth-microservice/internal/adapters/crypto"
	"auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handlers struct {
	Service *service.UserService
	JWT     *jwtutil.Manager
}

func NewHandlers(service *service.UserService, jwt *jwtutil.Manager) *Handlers {
	return &Handlers{
		Service: service,
		JWT:     jwt,
	}
}

type CheckUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req CheckUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}
	user, err := h.Service.UserByLogin(req.Login)
	if err != nil {
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !crypto.CheckPwd([]byte(req.Password), user.HashedPwd) {
		fmt.Fprintln(w, "Incorrect password ", http.StatusUnauthorized)
		return
	}

	token, err := h.JWT.Generate(user.Id, user.Role)
	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": token,
	})
}

type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error")
		log.Println(err)
		return
	}

	user, err := h.Service.AddUser(&domain.User{
		Login:     req.Login,
		HashedPwd: crypto.HashPassword(req.Password),
		Role:      domain.RoleBuyer,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	resp := map[string]any{
		"id":    user.Id,
		"login": user.Login,
		"role":  user.Role,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
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
				http.Error(w, "NO token", http.StatusUnauthorized)
				return
			}

			tokenStr := authHeader[len(prefix):]

			claims, err := jwtM.Parse(tokenStr)
			if err != nil {
				http.Error(w, "CANT PARSE", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "role", claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (h *Handlers) HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	usr, err := h.Service.UserByID(uint(userID))
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"user_id": userID,
		"role":    usr.Role,
	})
}
