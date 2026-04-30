package service_test

import (
	"auth-microservice/internal/core/domain"
	"auth-microservice/internal/core/service"
	"auth-microservice/internal/core/service/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type passwordCheckerStub struct {
	hashResult []byte
	checkOK    bool
}

func (s *passwordCheckerStub) CheckPassword(_ []byte, _ []byte) bool {
	return s.checkOK
}

func (s *passwordCheckerStub) HashPassword(_ string) []byte {
	return s.hashResult
}

func TestAddUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	checker := &passwordCheckerStub{hashResult: []byte("hashed-password")}
	userService := service.NewService(mockRepo, checker)

	expected := &domain.User{Id: 1, Login: "Test", Role: domain.RoleBuyer}
	mockRepo.On("AddUser", mock.MatchedBy(func(user *domain.User) bool {
		return user.Login == "Test" && user.Role == domain.RoleBuyer && string(user.HashedPwd) == "hashed-password"
	})).Return(expected, nil).Once()

	result, err := userService.AddUser("Test", "password123", domain.RoleBuyer)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestAddUserValidation(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	checker := &passwordCheckerStub{hashResult: []byte("hashed-password")}
	userService := service.NewService(mockRepo, checker)

	_, err := userService.AddUser("ab", "password123", domain.RoleBuyer)
	assert.ErrorIs(t, err, domain.ErrLoginTooShort)

	_, err = userService.AddUser("valid-login", "short", domain.RoleBuyer)
	assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
}

func TestUserByID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewService(mockRepo, &passwordCheckerStub{})
	input := int64(1)
	expected := &domain.User{Id: 1, Login: "Test", Role: domain.RoleBuyer}

	mockRepo.On("UserByID", input).Return(expected, nil).Once()

	result, err := userService.UserByID(input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUserByLogin(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewService(mockRepo, &passwordCheckerStub{checkOK: true})
	expected := &domain.User{Id: 1, Login: "Test", Role: domain.RoleBuyer, HashedPwd: []byte("hash")}

	mockRepo.On("UserByLogin", "Test").Return(expected, nil).Once()

	result, err := userService.UserByLogin("Test", "password123")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUserByLoginInvalidPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewService(mockRepo, &passwordCheckerStub{checkOK: false})
	storedUser := &domain.User{Id: 1, Login: "Test", Role: domain.RoleBuyer, HashedPwd: []byte("hash")}

	mockRepo.On("UserByLogin", "Test").Return(storedUser, nil).Once()

	result, err := userService.UserByLogin("Test", "bad-password")
	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
	mockRepo.AssertExpectations(t)
}

func TestUserByLoginRepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewService(mockRepo, &passwordCheckerStub{checkOK: true})
	repoErr := errors.New("repo failed")

	mockRepo.On("UserByLogin", "Test").Return(nil, repoErr).Once()

	result, err := userService.UserByLogin("Test", "password123")
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userService := service.NewService(mockRepo, &passwordCheckerStub{})
	input := int64(1)

	mockRepo.On("DeleteUser", input).Return(nil).Once()

	err := userService.DeleteUser(input)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
