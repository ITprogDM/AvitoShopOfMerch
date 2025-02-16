package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type StubUserRepository struct {
	GetUserByUsernameFunc func(username string) (*models.User, error)
	GetUserBalanceFunc    func(username string) (int, error)
	CreateUserFunc        func(user models.User) error
	UpdateUserBalanceFunc func(username string, newBalance int) error
	UserExistsFunc        func(username string) (bool, error)
}

func (s *StubUserRepository) GetUserByUsername(username string) (*models.User, error) {
	return s.GetUserByUsernameFunc(username)
}

func (s *StubUserRepository) GetUserBalance(username string) (int, error) {
	return s.GetUserBalanceFunc(username)
}

func (s *StubUserRepository) CreateUser(user models.User) error {
	return s.CreateUserFunc(user)
}

func (s *StubUserRepository) UpdateUserBalance(username string, newBalance int) error {
	return s.UpdateUserBalanceFunc(username, newBalance)
}

func (s *StubUserRepository) UserExists(username string) (bool, error) {
	return s.UserExistsFunc(username)
}

func TestAuthService_Login(t *testing.T) {
	// Хешируем пароль для теста
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// Создаем заглушку для UserRepository
	stubUserRepo := &StubUserRepository{
		GetUserByUsernameFunc: func(username string) (*models.User, error) {
			if username == "testuser" {
				return &models.User{
					Username: "testuser",
					Password: string(hashedPassword), // Используем корректный bcrypt-хеш
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	logger := logrus.New() // Инициализируем логгер
	authService := services.NewAuthService(stubUserRepo, "secret", logger)

	// Тест на успешный логин
	token, err := authService.Login("testuser", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Тест на неверный пароль
	_, err = authService.Login("testuser", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())

	// Тест на несуществующего пользователя
	_, err = authService.Login("nonexistent", "password")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestAuthService_Register(t *testing.T) {
	// Создаем заглушку для UserRepository
	stubUserRepo := &StubUserRepository{
		UserExistsFunc: func(username string) (bool, error) {
			if username == "existinguser" {
				return true, nil
			}
			return false, nil
		},
		CreateUserFunc: func(user models.User) error {
			return nil
		},
	}

	logger := logrus.New() // Инициализируем логгер
	authService := services.NewAuthService(stubUserRepo, "secret", logger)

	// Тест на успешную регистрацию
	token, err := authService.Register("newuser", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Тест на уже существующего пользователя
	_, err = authService.Register("existinguser", "password")
	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
}
