package handlers

import (
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserServiceForUsHand struct {
	mock.Mock
}

func (m *MockUserServiceForUsHand) GetBalance(username string) (int, error) {
	args := m.Called(username)
	return args.Int(0), args.Error(1)
}

func (m *MockUserServiceForUsHand) UserExists(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

// MockTransactionService - заглушка для TransactionService
type MockTransactionServiceForUsHand struct {
	mock.Mock
}

func (m *MockTransactionServiceForUsHand) GetReceivedTransactions(username string) ([]models.TransactionDetail, error) {
	args := m.Called(username)
	return args.Get(0).([]models.TransactionDetail), args.Error(1)
}

func (m *MockTransactionServiceForUsHand) GetSentTransactions(username string) ([]models.TransactionDetail, error) {
	args := m.Called(username)
	return args.Get(0).([]models.TransactionDetail), args.Error(1)
}

// MockInventoryService - заглушка для InventoryService
type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) GetInventory(username string) ([]models.InventoryItem, error) {
	args := m.Called(username)
	return args.Get(0).([]models.InventoryItem), args.Error(1)
}

func (m *MockInventoryService) AddToInventory(username string, itemType string, quantity int) error {
	args := m.Called(username, itemType, quantity)
	return args.Error(0)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) UserExists(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserService) GetBalance(username string) (int, error) {
	args := m.Called(username)
	return args.Int(0), args.Error(1)
}

// ================== ТЕСТ GetBalance ==================
func TestGetBalance(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		mockLogger := logrus.New()
		mockLogger.SetOutput(io.Discard)

		mockUserService := new(MockUserService) // Создаем мок сервиса
		mockUserService.On("GetBalance", mock.Anything).Return(0, errors.New("some error"))

		handler := handlers.NewUserHandler(mockUserService, nil, nil, mockLogger) // Создаем хэндлер с моками
		router.GET("/balance", func(c *gin.Context) {
			c.Set("username", "testuser") // Устанавливаем username в контекст
			handler.GetBalance(c)         // Вызываем хэндлер
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/balance", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		expected := map[string]interface{}{"error": "Failed to get balance"}
		var actual map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Equal(t, expected, actual)
	})
}

// ================== ТЕСТ GetUserInfo ==================
func TestGetUserInfo(t *testing.T) {
	t.Run("Balance_Fetch_Error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()

		mockLogger := logrus.New()
		mockLogger.SetOutput(io.Discard)

		mockUserService := new(MockUserService)
		mockTransactionService := new(MockTransactionService)
		mockInventoryService := new(MockInventoryService)

		// Мокаем ошибку при получении баланса
		mockUserService.On("GetBalance", "testuser").Return(0, errors.New("some error"))

		handler := handlers.NewUserHandler(mockUserService, mockTransactionService, mockInventoryService, mockLogger)
		router.GET("/info", func(c *gin.Context) {
			c.Set("username", "testuser")
			handler.GetUserInfo(c)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/info", nil)

		router.ServeHTTP(w, req)

		// Ожидаем 500 статус
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		expected := map[string]interface{}{"error": "Failed to fetch balance"}
		var actual map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// Проверяем, что API вернул правильное тело ошибки
		assert.Equal(t, expected, actual)
	})
}
