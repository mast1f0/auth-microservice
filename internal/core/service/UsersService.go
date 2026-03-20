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

func (service *UserService) UserCount() uint {
	return service.UsersRepository.UserCount()
}
func (service *UserService) UserByID(id uint) (*domain.User, error) {
	return service.UsersRepository.UserByID(id)
}

func (service *UserService) UserExists(login string) bool {
	return service.UsersRepository.UserExists(login)
}

func (service *UserService) AddUser(user *domain.User) (*domain.User, error) {
	return service.UsersRepository.AddUser(user)
}

func (service *UserService) UserByLogin(login string) (*domain.User, error) {
	return service.UsersRepository.UserByLogin(login)
}
