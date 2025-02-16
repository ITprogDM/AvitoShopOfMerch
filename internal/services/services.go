package services

import "ShopAvito/internal/models"

type PurchaseServiceInterface interface {
	BuyItem(username, itemName string, price int) error
	GetUserPurchases(username string) ([]models.Purchase, error)
}

type InventoryServiceInterface interface {
	GetInventory(username string) ([]models.InventoryItem, error)
	AddToInventory(username, itemType string, quantity int) error
}

type TransactionServiceInterface interface {
	TransferCoins(fromUser, toUser string, amount int) error
	GetReceivedTransactions(username string) ([]models.TransactionDetail, error)
	GetSentTransactions(username string) ([]models.TransactionDetail, error)
}

type UserServiceInterface interface {
	UserExists(username string) (bool, error)
	GetBalance(username string) (int, error)
}

type AuthServiceInterface interface {
	GenerateToken(username string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	Login(username, password string) (string, error)
	Register(username, password string) (string, error)
}
