package models

import (
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Пароль не возвращается в JSON
	Balance  int    `json:"balance"`
}

// Transaction - структура для перевода монет
type Transaction struct {
	ID       int       `json:"id"`
	FromUser string    `json:"from_user"`
	ToUser   string    `json:"to_user"`
	Amount   int       `json:"amount"`
	Time     time.Time `json:"timestamp"`
}

// Purchase - структура для покупки товара
type Purchase struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	ItemName string    `json:"item_name"`
	Price    int       `json:"price"`
	Time     time.Time `json:"timestamp"`
}

// AuthRequest - запрос на авторизацию
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse - ответ с JWT токеном
type AuthResponse struct {
	Token string `json:"token"`
}

// SendCoinRequest - запрос на перевод монет
type SendCoinRequest struct {
	ToUser string `json:"to_user"`
	Amount int    `json:"amount"`
}

// InfoResponse - ответ с информацией о пользователе
type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

// InventoryItem - предмет в инвентаре
type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

// CoinHistory - история операций с монетами
type CoinHistory struct {
	Received []TransactionDetail `json:"received"`
	Sent     []TransactionDetail `json:"sent"`
}

// TransactionDetail - детали транзакции (для истории)
type TransactionDetail struct {
	FromUser string `json:"from_user,omitempty"`
	ToUser   string `json:"to_user,omitempty"`
	Amount   int    `json:"amount"`
}
