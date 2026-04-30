package http

import (
	"auth-microservice/internal/adapters/crypto"
	"auth-microservice/internal/adapters/http/handlers"
	jwtutil "auth-microservice/internal/adapters/jwt"
	"auth-microservice/internal/core/service"
	"auth-microservice/internal/core/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterExposesPublicAndProtectedRoutes(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	h := handlers.NewHandlers(service.NewService(mockRepo, crypto.NewBcryptHasher()), &jwtutil.Manager{Secret: []byte("secret")})
	router := NewRouter(h)

	registerReq := httptest.NewRequest(http.MethodPost, "/register", nil)
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	if registerResp.Code == http.StatusNotFound {
		t.Fatal("expected /register route to exist")
	}

	profileReq := httptest.NewRequest(http.MethodGet, "/profile", nil)
	profileResp := httptest.NewRecorder()
	router.ServeHTTP(profileResp, profileReq)
	if profileResp.Code != http.StatusUnauthorized {
		t.Fatalf("expected /profile to be protected, got %d", profileResp.Code)
	}
}
