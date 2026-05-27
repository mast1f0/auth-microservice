package ports

import "errors"

var (
	ErrFailedToLoad = errors.New("failed to load")
	ErrUserNotFound = errors.New("user not found")
)
