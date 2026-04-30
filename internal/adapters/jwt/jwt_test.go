package jwtutil

import (
	"auth-microservice/internal/core/domain"
	"testing"
)

func TestManager_GenerateAndParse(t *testing.T) {
	manager := &Manager{Secret: []byte("secret")}

	token, err := manager.Generate(42, domain.RoleSeller)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("expected user id 42, got %d", claims.UserID)
	}
	if claims.Role != domain.RoleSeller {
		t.Fatalf("expected role %q, got %q", domain.RoleSeller, claims.Role)
	}
}

func TestManager_ParseInvalidToken(t *testing.T) {
	manager := &Manager{Secret: []byte("secret")}

	if _, err := manager.Parse("bad-token"); err == nil {
		t.Fatal("expected parse error for invalid token")
	}
}
