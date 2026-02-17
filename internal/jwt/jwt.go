package jwtutil

import "github.com/golang-jwt/jwt/v5"

type Manager struct {
	Secret []byte
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
