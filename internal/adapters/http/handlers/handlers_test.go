package handlers_test

import (
	"auth-microservice/internal/adapters/http/handlers"
	"auth-microservice/internal/core/domain"
	service2 "auth-microservice/internal/core/service"
	"auth-microservice/internal/core/service/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRegister(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	mockService := service2.NewService(mockRepo)

	h := handlers.Handlers{Service: mockService}
	bodyReq := handlers.RegisterUser{
		Login:    "login",
		Password: "password123",
	}
	body, _ := json.Marshal(bodyReq)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	expectedUser := &domain.User{
		Id:    1,
		Login: "login",
		Role:  domain.RoleBuyer,
	}

	mockRepo.On("AddUser", mock.MatchedBy(func(user *domain.User) bool {
		return user.Login == "login" &&
			user.Role == domain.RoleBuyer &&
			user.HashedPwd != nil &&
			len(user.HashedPwd) > 0
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
	assert.Equal(t, "buyer", response["role"])

	mockRepo.AssertExpectations(t)
}
