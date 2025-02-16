package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type StubInventoryRepository struct {
	GetInventoryFunc   func(username string) ([]models.InventoryItem, error)
	AddToInventoryFunc func(username, itemType string, quantity int) error
}

func (s *StubInventoryRepository) GetInventory(username string) ([]models.InventoryItem, error) {
	return s.GetInventoryFunc(username)
}

func (s *StubInventoryRepository) AddToInventory(username, itemType string, quantity int) error {
	return s.AddToInventoryFunc(username, itemType, quantity)
}

func TestInventoryService_GetInventory(t *testing.T) {
	// Создаем заглушку для InventoryRepository
	stubInventoryRepo := &StubInventoryRepository{
		GetInventoryFunc: func(username string) ([]models.InventoryItem, error) {
			if username == "testuser" {
				return []models.InventoryItem{
					{Type: "t-shirt", Quantity: 2},
					{Type: "cup", Quantity: 1},
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	// Создаем InventoryService с заглушкой
	inventoryService := services.NewInventoryService(stubInventoryRepo)

	// Тест на успешное получение инвентаря
	items, err := inventoryService.GetInventory("testuser")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "t-shirt", items[0].Type)
	assert.Equal(t, 2, items[0].Quantity)
	assert.Equal(t, "cup", items[1].Type)
	assert.Equal(t, 1, items[1].Quantity)

	// Тест на несуществующего пользователя
	_, err = inventoryService.GetInventory("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestInventoryService_AddToInventory(t *testing.T) {
	// Создаем заглушку для InventoryRepository
	stubInventoryRepo := &StubInventoryRepository{
		AddToInventoryFunc: func(username, itemType string, quantity int) error {
			if username == "testuser" && itemType == "t-shirt" {
				return nil
			}
			return errors.New("failed to add item")
		},
	}

	// Создаем InventoryService с заглушкой
	inventoryService := services.NewInventoryService(stubInventoryRepo)

	// Тест на успешное добавление предмета в инвентарь
	err := inventoryService.AddToInventory("testuser", "t-shirt", 1)
	assert.NoError(t, err)

	// Тест на ошибку при добавлении предмета
	err = inventoryService.AddToInventory("testuser", "invalid-item", 1)
	assert.Error(t, err)
	assert.Equal(t, "failed to add item", err.Error())
}
