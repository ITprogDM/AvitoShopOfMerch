package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
)

type InventoryService struct {
	inventoryRepo repository.InventoryRepositoryInterface
}

func NewInventoryService(inventoryRepo repository.InventoryRepositoryInterface) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (s *InventoryService) GetInventory(username string) ([]models.InventoryItem, error) {
	return s.inventoryRepo.GetInventory(username)
}

func (s *InventoryService) AddToInventory(username, itemType string, quantity int) error {
	return s.inventoryRepo.AddToInventory(username, itemType, quantity)
}
