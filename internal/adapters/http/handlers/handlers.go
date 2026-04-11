package handlers

import (
	"auth-microservice/internal/adapters/crypto"
	"auth-microservice/internal/adapters/http/handlers/helpers"
	"auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/service"
	"context"
	"encoding/json"
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
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.Service.UserByLogin(req.Login)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !crypto.CheckPwd([]byte(req.Password), user.HashedPwd) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "incorrect password")
		return
	}

	token, err := h.JWT.Generate(user.Id, user.Role)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"access_token": token})
}

type RegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Login) == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "login must be at least symbols")
		return
	}

	if len(req.Password) < 8 {
		helpers.RespondWithError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}
	user, err := h.Service.AddUser(&domain.User{
		Login:     req.Login,
		HashedPwd: crypto.HashPassword(req.Password),
		Role:      domain.RoleBuyer,
	})
	if err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{
		"id":         user.Id,
		"login":      user.Login,
		"created_at": user.CreatedAt,
	})
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
		helpers.RespondWithError(w, http.StatusInternalServerError, "CANT_GET_USER_ID")
		return
	}
	usr, err := h.Service.UserByID(uint(userID))
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, map[string]any{
		"user_id": usr.Id,
		"role":    usr.Role,
	})
}

func (h *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		helpers.RespondWithError(w, http.StatusInternalServerError, "CANT_GET_USER_ID")
		return
	}
	err := h.Service.DeleteUser(uint(userID))
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusNoContent, nil)
}
