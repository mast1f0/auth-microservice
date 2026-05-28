package handlers

import (
	"auth-microservice/internal/adapters/http/handlers/dto"
	"auth-microservice/internal/adapters/http/handlers/helpers"
	jwtutil "auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/service"
	"encoding/json"
	"errors"
	"net/http"
)

const defaultRole = domain.RoleBuyer

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

func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.Service.UserByLogin(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrUserNotFound) {
			helpers.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := h.JWT.Generate(user.Id, user.Role)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"access_token": token})
}

func (h *Handlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.Service.AddUser(req.Login, req.Password, defaultRole)
	if err != nil {
		if errors.Is(err, service.ErrLoginTooShort) || errors.Is(err, service.ErrPasswordTooShort) {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusCreated, map[string]any{
		"id":         user.Id,
		"login":      user.Login,
		"created_at": user.CreatedAt,
	})
}

func (h *Handlers) HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		helpers.RespondWithError(w, http.StatusInternalServerError, "CANT_GET_USER_ID")
		return
	}
	usr, err := h.Service.UserByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			helpers.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
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
	userRole, ok := r.Context().Value("role").(domain.Role)
	if !ok {
		helpers.RespondWithError(w, http.StatusInternalServerError, "CANT_GET_USER_ROLE")
		return
	}

	var req dto.DeleteUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = req.Validate(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if userRole != domain.RoleAdmin && userID != req.UserID {
		helpers.RespondWithError(w, http.StatusForbidden, "Permission denied")
		return
	}
	err = h.Service.DeleteUser(req.UserID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.AllUsers()
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, users)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		helpers.RespondWithError(w, http.StatusInternalServerError, "CANT_GET_USER_ID")
		return
	}
	userRole, ok := r.Context().Value("role").(domain.Role)
	if !ok {
		helpers.RespondWithError(w, http.StatusInternalServerError, "CANT_GET_USER_ROLE")
		return
	}
	if userRole != domain.RoleAdmin {
		helpers.RespondWithError(w, http.StatusForbidden, "Permission denied")
		return
	}

	var req dto.UpdateUserDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := req.Valid(); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = h.Service.UpdateUser(userID, domain.Role(req.Role))
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, "updated")

}
