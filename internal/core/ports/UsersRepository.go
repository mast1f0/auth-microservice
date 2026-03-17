package ports

import (
	"auth-microservice/internal/core/domain"
)

type UsersRepository interface {
	UserCount() uint
	UserByID(id uint) (*domain.User, error)
	UserExists(login string) bool
	AddUser(user *domain.User) error
	UserByLogin(login string) (*domain.User, error)
}
