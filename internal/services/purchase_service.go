package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
	"errors"
	"github.com/sirupsen/logrus"
)

type PurchaseService struct {
	purchaseRepo  repository.PurchaseRepositoryInterface
	userRepo      repository.UserRepositoryInterface
	inventoryRepo repository.InventoryRepositoryInterface
	log           *logrus.Logger
}

func NewPurchaseService(purchaseRepo repository.PurchaseRepositoryInterface, userRepo repository.UserRepositoryInterface, inventoryRepo repository.InventoryRepositoryInterface, log *logrus.Logger) *PurchaseService {
	return &PurchaseService{
		purchaseRepo:  purchaseRepo,
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
		log:           log,
	}
}

// Покупка товара
func (s *PurchaseService) BuyItem(username, itemName string, price int) error {
	// Проверяем баланс
	balance, err := s.userRepo.GetUserBalance(username)
	if err != nil {
		s.log.Errorf("Error getting user balance: %v", err)
		return err
	}
	if balance < price {
		return errors.New("insufficient funds")
	}

	// Покупаем предмет и обновляем инвентарь
	err = s.purchaseRepo.BuyItem(username, itemName, price)
	if err != nil {
		s.log.Errorf("Error buying item: %v", err)
	}

	s.log.Infof("Adding item to inventory: username=%s, item=%s", username, itemName)
	err = s.inventoryRepo.AddToInventory(username, itemName, 0)
	if err != nil {
		s.log.Errorf("Error adding item to inventory: %v", err)
		return err
	}

	return err
}

// Получение списка купленных товаров
func (s *PurchaseService) GetUserPurchases(username string) ([]models.Purchase, error) {
	return s.purchaseRepo.GetUserPurchases(username)
}
