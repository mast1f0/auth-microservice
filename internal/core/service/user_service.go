package service

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/ports"
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
	return service.UsersRepository.UserByID(id)
}

func (service *UserService) AddUser(login string, password string, role domain.Role) (*domain.User, error) {
	if len(login) < 3 {
		return nil, domain.ErrLoginTooShort
	}

	if len(password) < 8 {
		return nil, domain.ErrPasswordTooShort
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
		return nil, err
	}
	if service.PasswordChecker.CheckPassword([]byte(password), usr.HashedPwd) {
		return usr, nil
	}
	return nil, domain.ErrInvalidCredentials
}

func (service *UserService) DeleteUser(id int64) error {
	return service.UsersRepository.DeleteUser(id)
}
