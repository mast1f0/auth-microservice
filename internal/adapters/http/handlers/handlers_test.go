package handlers_test

import (
	"auth-microservice/internal/adapters/crypto"
	httpadapter "auth-microservice/internal/adapters/http"
	"auth-microservice/internal/adapters/http/handlers"
	"auth-microservice/internal/adapters/http/handlers/dto"
	jwtutil "auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/service"
	"auth-microservice/internal/core/service/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestHandler(mockRepo *mocks.MockUserRepository) *handlers.Handlers {
	mockService := service.NewService(mockRepo, crypto.NewBcryptHasher())
	jwtManager := &jwtutil.Manager{Secret: []byte("secret")}
	return handlers.NewHandlers(mockService, jwtManager)
}

func TestHandleRegister_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)

	bodyReq := dto.RegisterUser{Login: "login", Password: "password123"}
	body, _ := json.Marshal(bodyReq)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	expectedUser := &domain.User{Id: 1, Login: "login", Role: domain.RoleBuyer}
	mockRepo.On("AddUser", mock.MatchedBy(func(user *domain.User) bool {
		return user.Login == "login" && user.Role == domain.RoleBuyer && len(user.HashedPwd) > 0
	})).Return(expectedUser, nil).Once()

	h.HandleRegister(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "login", response["login"])

	mockRepo.AssertExpectations(t)
}

func TestHandleRegister_ValidationError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)

	body, _ := json.Marshal(dto.RegisterUser{Login: "ab", Password: "password123"})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.HandleRegister(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockRepo.AssertNotCalled(t, "AddUser", mock.Anything)
}

func TestHandleLogin_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)

	hashedPwd := crypto.NewBcryptHasher().HashPassword("password123")
	expectedUser := &domain.User{Id: 1, Login: "login", HashedPwd: hashedPwd, Role: domain.RoleBuyer}
	mockRepo.On("UserByLogin", "login").Return(expectedUser, nil).Once()

	bodyReq := dto.LoginReq{Login: "login", Password: "password123"}
	body, _ := json.Marshal(bodyReq)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.HandleLogin(rr, req)

	resp := rr.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["access_token"])

	mockRepo.AssertExpectations(t)
}

func TestHandleLogin_InvalidCredentials(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)

	mockRepo.On("UserByLogin", "login").Return(nil, service.ErrUserNotFound).Once()

	body, _ := json.Marshal(dto.LoginReq{Login: "login", Password: "password123"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.HandleLogin(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleProfile_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)
	userID := int64(7)

	mockRepo.On("UserByID", userID).Return(&domain.User{Id: userID, Role: domain.RoleSeller}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	rr := httptest.NewRecorder()

	h.HandleProfile(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandleProfile_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)
	userID := int64(7)

	mockRepo.On("UserByID", userID).Return(nil, service.ErrUserNotFound).Once()

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	rr := httptest.NewRecorder()

	h.HandleProfile(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)
	userID := int64(3)

	mockRepo.On("DeleteUser", userID).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", userID))
	rr := httptest.NewRecorder()

	h.DeleteUser(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body.String())
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware(t *testing.T) {
	jwtManager := &jwtutil.Manager{Secret: []byte("secret")}
	token, err := jwtManager.Generate(9, domain.RoleAdmin)
	assert.NoError(t, err)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, int64(9), r.Context().Value("user_id"))
		assert.Equal(t, domain.RoleAdmin, r.Context().Value("role"))
		w.WriteHeader(http.StatusOK)
	})

	middleware := handlers.AuthMiddleware(jwtManager)
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	middleware(next).ServeHTTP(rr, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRouter_WiresProtectedEndpoints(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)
	router := httpadapter.NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestHandleLogin_InternalError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := newTestHandler(mockRepo)
	repoErr := errors.New("db down")
	mockRepo.On("UserByLogin", "login").Return(nil, repoErr).Once()

	body, _ := json.Marshal(dto.LoginReq{Login: "login", Password: "password123"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.HandleLogin(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
