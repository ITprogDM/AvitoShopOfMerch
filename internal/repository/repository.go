package repository

import "ShopAvito/internal/models"

type PurchaseRepositoryInterface interface {
	BuyItem(username, itemName string, price int) error
	GetUserPurchases(username string) ([]models.Purchase, error)
}

type InventoryRepositoryInterface interface {
	GetInventory(username string) ([]models.InventoryItem, error)
	AddToInventory(username, itemType string, quantity int) error
}

type TransactionRepositoryInterface interface {
	GetUserID(username string) (int, error)
	TransferCoins(fromUser, toUser string, amount int) error
	GetReceivedTransactions(username string) ([]models.Transaction, error)
	GetSentTransactions(username string) ([]models.Transaction, error)
}

type UserRepositoryInterface interface {
	CreateUser(user models.User) error
	GetUserBalance(username string) (int, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserBalance(username string, newBalance int) error
	UserExists(username string) (bool, error)
}
