package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type StubUserRepositoryForUser struct {
	CreateUserFunc        func(user models.User) error
	GetUserBalanceFunc    func(username string) (int, error)
	GetUserByUsernameFunc func(username string) (*models.User, error)
	UpdateUserBalanceFunc func(username string, newBalance int) error
	UserExistsFunc        func(username string) (bool, error)
}

func (s *StubUserRepositoryForUser) CreateUser(user models.User) error {
	return s.CreateUserFunc(user)
}

func (s *StubUserRepositoryForUser) GetUserBalance(username string) (int, error) {
	return s.GetUserBalanceFunc(username)
}

func (s *StubUserRepositoryForUser) GetUserByUsername(username string) (*models.User, error) {
	return s.GetUserByUsernameFunc(username)
}

func (s *StubUserRepositoryForUser) UpdateUserBalance(username string, newBalance int) error {
	return s.UpdateUserBalanceFunc(username, newBalance)
}

func (s *StubUserRepositoryForUser) UserExists(username string) (bool, error) {
	return s.UserExistsFunc(username)
}

func TestUserService_UserExists(t *testing.T) {
	// Создаем заглушку для UserRepository
	stubUserRepo := &StubUserRepositoryForUser{
		UserExistsFunc: func(username string) (bool, error) {
			if username == "existinguser" {
				return true, nil
			}
			return false, nil
		},
	}

	// Инициализируем логгер с отключенным выводом
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Отключаем вывод логов в консоль

	// Создаем UserService с заглушкой и логгером
	userService := services.NewUserService(stubUserRepo, logger)

	// Тест на существующего пользователя
	exists, err := userService.UserExists("existinguser")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Тест на несуществующего пользователя
	exists, err = userService.UserExists("nonexistent")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Тест на ошибку в репозитории
	stubUserRepo.UserExistsFunc = func(username string) (bool, error) {
		return false, errors.New("database error")
	}
	_, err = userService.UserExists("erroruser")
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestUserService_GetBalance(t *testing.T) {
	// Создаем заглушку для UserRepository
	stubUserRepo := &StubUserRepository{
		GetUserBalanceFunc: func(username string) (int, error) {
			if username == "richuser" {
				return 1000, nil
			}
			return 0, errors.New("user not found")
		},
	}

	// Инициализируем логгер с отключенным выводом
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Отключаем вывод логов в консоль

	// Создаем UserService с заглушкой и логгером
	userService := services.NewUserService(stubUserRepo, logger)

	// Тест на успешное получение баланса
	balance, err := userService.GetBalance("richuser")
	assert.NoError(t, err)
	assert.Equal(t, 1000, balance)

	// Тест на несуществующего пользователя
	_, err = userService.GetBalance("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}
