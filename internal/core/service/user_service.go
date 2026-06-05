package service

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/ports"
	"context"
	"errors"
)

type UserService struct {
	UsersRepository ports.UsersRepository
	PasswordChecker ports.PasswordCheckerRepository
}

func NewService(usersRepository ports.UsersRepository, passwordCheckerRepository ports.PasswordCheckerRepository) *UserService {
	return &UserService{
		UsersRepository: usersRepository,
		PasswordChecker: passwordCheckerRepository,
	}
}

func (service *UserService) UserByID(ctx context.Context, id int64) (*domain.User, error) {
	if id < 1 {
		return nil, ErrInvalidId
	}
	return service.UsersRepository.UserByID(ctx, id)
}

func (service *UserService) AddUser(ctx context.Context, login string, password string, role domain.Role) (*domain.User, error) {
	if len(login) < 3 {
		return nil, ErrLoginTooShort
	}

	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	return service.UsersRepository.AddUser(ctx, &domain.User{
		Login:     login,
		HashedPwd: service.PasswordChecker.HashPassword(password),
		Role:      role,
	})
}

func (service *UserService) UserByLogin(ctx context.Context, login string, password string) (*domain.User, error) {
	usr, err := service.UsersRepository.UserByLogin(ctx, login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if service.PasswordChecker.CheckPassword([]byte(password), usr.HashedPwd) {
		return usr, nil
	}
	return nil, ErrInvalidCredentials
}

func (service *UserService) DeleteUser(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrInvalidId
	}
	return service.UsersRepository.DeleteUser(ctx, id)
}

func (service *UserService) AllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := service.UsersRepository.GetAllUsers(ctx)
	if err != nil {
		if errors.Is(err, ports.ErrFailedToLoad) {
			return []*domain.User{}, ErrInternalServer
		}
		return []*domain.User{}, err
	}
	return users, nil
}

func (service *UserService) UpdateUser(ctx context.Context, userId int64, role domain.Role) error {
	if userId < 1 {
		return ErrInvalidId
	}
	err := service.UsersRepository.UpdateUser(ctx, userId, role)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}
