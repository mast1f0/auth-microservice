package service

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/ports"
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

func (service *UserService) UserByID(id int64) (*domain.User, error) {
	if id < 1 {
		return nil, ErrInvalidId
	}
	return service.UsersRepository.UserByID(id)
}

func (service *UserService) AddUser(login string, password string, role domain.Role) (*domain.User, error) {
	if len(login) < 3 {
		return nil, ErrLoginTooShort
	}

	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	return service.UsersRepository.AddUser(&domain.User{
		Login:     login,
		HashedPwd: service.PasswordChecker.HashPassword(password),
		Role:      role,
	})
}

func (service *UserService) UserByLogin(login string, password string) (*domain.User, error) {
	usr, err := service.UsersRepository.UserByLogin(login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if service.PasswordChecker.CheckPassword([]byte(password), usr.HashedPwd) {
		return usr, nil
	}
	return nil, ErrInvalidCredentials
}

func (service *UserService) DeleteUser(id int64) error {
	if id < 1 {
		return ErrInvalidId
	}
	return service.UsersRepository.DeleteUser(id)
}

func (service *UserService) AllUsers() ([]*domain.User, error) {
	users, err := service.UsersRepository.GetAllUsers()
	if err != nil {
		if errors.Is(err, ports.ErrFailedToLoad) {
			return []*domain.User{}, ErrInternalServer
		}
		return []*domain.User{}, err
	}
	return users, nil
}

func (service *UserService) UpdateUser(userId int64, role domain.Role) error {
	if userId < 1 {
		return ErrInvalidId
	}
	err := service.UsersRepository.UpdateUser(userId, role)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}
