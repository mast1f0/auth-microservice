package ports

import (
	"auth-microservice/internal/core/domain"
)

type UsersRepository interface {
	UserByID(id int64) (*domain.User, error)
	AddUser(user *domain.User) (*domain.User, error)
	UserByLogin(login string) (*domain.User, error)
	DeleteUser(id int64) error
	GetAllUsers() ([]*domain.User, error)
	UpdateUser(userId int64, newRole domain.Role) error
}
