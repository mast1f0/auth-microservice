package jwtutil

import (
	"auth-microservice/internal/core/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	Secret []byte
}

type Claims struct {
	UserID int64       `json:"user_id"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func (m *Manager) Generate(userID int64, role domain.Role) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.Secret)
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return m.Secret, nil
	},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
