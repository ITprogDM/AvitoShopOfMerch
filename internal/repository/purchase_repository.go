package repository

import (
	"ShopAvito/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PurchaseRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewPurchaseRepository(db *pgxpool.Pool, log *logrus.Logger) *PurchaseRepository {
	return &PurchaseRepository{
		db:  db,
		log: log,
	}
}

func (r *PurchaseRepository) BuyItem(username, itemName string, price int) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		r.log.Errorf("Error starting transaction: %s", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(context.Background())
			r.log.Errorf("Transaction panicked: %v", p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil {
				r.log.Errorf("Transaction rollback failed: %v", rollbackErr)
			}
		}
	}()

	// Проверяем, существует ли пользователь
	var userID int
	err = tx.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		r.log.Errorf("Failed to get user ID for username %s: %v", username, err)
		return err
	}

	// Проверяем, достаточно ли баланса
	var currentBalance int
	err = tx.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
	if err != nil {
		r.log.Errorf("Failed to get balance for user %s: %v", username, err)
		return err
	}
	if currentBalance < price {
		r.log.Errorf("Insufficient balance for user %s: %d < %d", username, currentBalance, price)
		return fmt.Errorf("insufficient balance")
	}
	// Вычитаем баланс
	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE id = $2", price, userID)
	if err != nil {
		r.log.Errorf("Failed to update balance for user %s: %v", username, err)
		return err
	}

	// Добавляем запись в purchases
	_, err = tx.Exec(context.Background(),
		"INSERT INTO purchases (user_id, item_name, price) VALUES ($1, $2, $3)",
		userID, itemName, price)
	if err != nil {
		r.log.Errorf("Failed to insert purchase record for user %s: %v", username, err)
		return err
	}

	// Добавляем в инвентарь или увеличиваем количество
	_, err = tx.Exec(context.Background(),
		`INSERT INTO inventory (user_id, item_type, quantity) 
         VALUES ($1, $2, 1) 
         ON CONFLICT (user_id, item_type) 
         DO UPDATE SET quantity = inventory.quantity + 1`,
		userID, itemName)
	if err != nil {
		r.log.Errorf("Failed to update inventory for user %s: %v", username, err)
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		r.log.Errorf("Failed to commit transaction for user %s: %v", username, err)
		return err
	}
	r.log.Infof("Purchase successful for user %s: %s", username, itemName)
	return nil
}

func (r *PurchaseRepository) GetUserPurchases(username string) ([]models.Purchase, error) {
	var userID int
	err := r.db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		r.log.Errorf("Failed to get user ID for username %s: %v", username, err)
		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		"SELECT item_name, price, timestamp FROM purchases WHERE user_id = $1", userID)
	if err != nil {
		r.log.Errorf("Failed to fetch purchases for user %s: %v", username, err)
		return nil, err
	}
	defer rows.Close()

	var purchases []models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		err = rows.Scan(&purchase.ItemName, &purchase.Price, &purchase.Time)
		if err != nil {
			r.log.Errorf("Failed to scan purchase for user %s: %v", username, err)
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	if err = rows.Err(); err != nil {
		r.log.Errorf("Error iterating over purchases for user %s: %v", username, err)
		return nil, err
	}

	r.log.Infof("Fetched %d purchases for user %s", len(purchases), username)
	return purchases, nil
}
