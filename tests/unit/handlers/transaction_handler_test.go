package handlers

import (
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) TransferCoins(fromUser, toUser string, amount int) error {
	args := m.Called(fromUser, toUser, amount)
	return args.Error(0)
}

func (m *MockTransactionService) GetReceivedTransactions(username string) ([]models.TransactionDetail, error) {
	args := m.Called(username)
	return args.Get(0).([]models.TransactionDetail), args.Error(1)
}

func (m *MockTransactionService) GetSentTransactions(username string) ([]models.TransactionDetail, error) {
	args := m.Called(username)
	return args.Get(0).([]models.TransactionDetail), args.Error(1)
}

// **Тестируем SendCoins**
func TestSendCoins_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := logrus.New()
	mockLogger.SetOutput(io.Discard)

	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService, mockLogger)

	router := gin.Default()
	router.POST("/send", handler.SendCoins)

	requestBody, _ := json.Marshal(models.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	})

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockService.On("TransferCoins", "testuser", "receiver", 100).Return(nil)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set("username", "testuser")

	handler.SendCoins(ctx)

	require.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// **Тестируем ошибку**
func TestSendCoins_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockLogger := logrus.New()
	mockLogger.SetOutput(io.Discard)

	mockService := new(MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService, mockLogger)

	router := gin.Default()
	router.POST("/send", handler.SendCoins)

	requestBody, _ := json.Marshal(models.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	})

	req, _ := http.NewRequest(http.MethodPost, "/send", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockService.On("TransferCoins", "testuser", "receiver", 100).Return(errors.New("transfer failed"))

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set("username", "testuser")

	handler.SendCoins(ctx)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
