package service_test

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/service"
	"auth-microservice/internal/core/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {

	mockRepo := new(mocks.MockUserRepository)

	service := service.NewService(mockRepo)

	input := domain.User{
		Login:     "Test",
		HashedPwd: []byte("hash"),
		Role:      domain.RoleBuyer,
	}

	expected := &domain.User{
		Id:    1,
		Login: "Test",
		Role:  domain.RoleBuyer,
	}

	mockRepo.On("AddUser", &input).Return(expected, nil)

	result, err := service.AddUser(&input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestUserByID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := service.NewService(mockRepo)
	input := uint(1)
	expected := &domain.User{
		Id:    1,
		Login: "Test",
		Role:  domain.RoleBuyer,
	}

	mockRepo.On("UserByID", input).Return(expected, nil)

	result, err := service.UserByID(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUserByLogin(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := service.NewService(mockRepo)
	input := "Test"
	expected := &domain.User{
		Id:    1,
		Login: input,
		Role:  domain.RoleBuyer,
	}
	mockRepo.On("UserByLogin", input).Return(expected, nil)
	result, err := service.UserByLogin(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := service.NewService(mockRepo)

	input := uint(1)

	mockRepo.On("DeleteUser", input).Return(nil)

	err := service.DeleteUser(input)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
