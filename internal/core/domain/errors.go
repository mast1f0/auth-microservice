package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrLoginTooShort      = errors.New("login must be at least 3 characters long")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
)
