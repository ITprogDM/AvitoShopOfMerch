package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type StubPurchaseRepository struct {
	BuyItemFunc          func(username, itemName string, price int) error
	GetUserPurchasesFunc func(username string) ([]models.Purchase, error)
}

func (s *StubPurchaseRepository) BuyItem(username, itemName string, price int) error {
	return s.BuyItemFunc(username, itemName, price)
}

func (s *StubPurchaseRepository) GetUserPurchases(username string) ([]models.Purchase, error) {
	return s.GetUserPurchasesFunc(username)
}

func TestPurchaseService_GetUserPurchases(t *testing.T) {
	// Создаем заглушку для PurchaseRepository
	stubPurchaseRepo := &StubPurchaseRepository{
		GetUserPurchasesFunc: func(username string) ([]models.Purchase, error) {
			if username == "testuser" {
				return []models.Purchase{
					{
						ItemName: "t-shirt",
						Price:    80,
						Time:     time.Now(),
					},
					{
						ItemName: "cup",
						Price:    20,
						Time:     time.Now(),
					},
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	// Инициализируем логгер
	logger := logrus.New()

	// Создаем PurchaseService с заглушкой и логгером
	purchaseService := services.NewPurchaseService(stubPurchaseRepo, nil, nil, logger)

	// Тест на успешное получение списка покупок
	purchases, err := purchaseService.GetUserPurchases("testuser")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(purchases))
	assert.Equal(t, "t-shirt", purchases[0].ItemName)
	assert.Equal(t, 80, purchases[0].Price)
	assert.Equal(t, "cup", purchases[1].ItemName)
	assert.Equal(t, 20, purchases[1].Price)

	// Тест на несуществующего пользователя
	_, err = purchaseService.GetUserPurchases("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}
