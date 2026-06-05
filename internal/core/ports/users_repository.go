package ports

import (
	"auth-microservice/internal/core/domain"
	"context"
)

type UsersRepository interface {
	UserByID(ctx context.Context, id int64) (*domain.User, error)
	AddUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UserByLogin(ctx context.Context, login string) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUser(ctx context.Context, userId int64, newRole domain.Role) error
}
