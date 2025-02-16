package handlers

import (
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockPurchaseService struct {
	mock.Mock
}

func (m *MockPurchaseService) BuyItem(username, item string, price int) error {
	args := m.Called(username, item, price)
	return args.Error(0)
}

func (m *MockPurchaseService) GetUserPurchases(username string) ([]models.Purchase, error) {
	args := m.Called(username)
	return args.Get(0).([]models.Purchase), args.Error(1)
}

func TestBuyItem_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPurchaseService)
	handler := handlers.NewPurchaseHandler(mockService, nil, logrus.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "item", Value: "t-shirt"}}
	c.Set("username", "testuser")

	mockService.On("BuyItem", "testuser", "t-shirt", 80).Return(nil)

	handler.BuyItem(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Purchase successful"}`, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestBuyItem_InvalidItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPurchaseService)
	handler := handlers.NewPurchaseHandler(mockService, nil, logrus.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "item", Value: "unknown"}}
	c.Set("username", "testuser")

	handler.BuyItem(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "Invalid item"}`, w.Body.String())
}

func TestBuyItem_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPurchaseService)
	handler := handlers.NewPurchaseHandler(mockService, nil, logrus.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "item", Value: "t-shirt"}}
	c.Set("username", "testuser")

	mockService.On("BuyItem", "testuser", "t-shirt", 80).Return(assert.AnError)

	handler.BuyItem(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
