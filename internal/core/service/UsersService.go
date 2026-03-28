package service

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/ports"
)

type UserService struct {
	UsersRepository ports.UsersRepository
}

func NewService(usersRepository ports.UsersRepository) *UserService {
	return &UserService{
		UsersRepository: usersRepository,
	}
}

func (service *UserService) UserByID(id uint) (*domain.User, error) {
	return service.UsersRepository.UserByID(id)
}

func (service *UserService) AddUser(user *domain.User) (*domain.User, error) {
	return service.UsersRepository.AddUser(user)
}

func (service *UserService) UserByLogin(login string) (*domain.User, error) {
	return service.UsersRepository.UserByLogin(login)
}

func (service *UserService) DeleteUser(id uint) error {
	return service.UsersRepository.DeleteUser(id)
}
