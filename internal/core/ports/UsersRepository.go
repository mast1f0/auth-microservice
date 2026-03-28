package ports

import (
	"auth-microservice/internal/core/domain"
)

type UsersRepository interface {
	UserByID(id uint) (*domain.User, error)
	AddUser(user *domain.User) (*domain.User, error)
	UserByLogin(login string) (*domain.User, error)
	DeleteUser(id uint) error
}
