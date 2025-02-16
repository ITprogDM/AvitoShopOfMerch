package repository

import (
	"ShopAvito/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type InventoryRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewInventoryRepository(db *pgxpool.Pool, log *logrus.Logger) *InventoryRepository {
	return &InventoryRepository{
		db:  db,
		log: log,
	}
}

func (r *InventoryRepository) GetInventory(username string) ([]models.InventoryItem, error) {
	var userID int
	err := r.db.QueryRow(context.Background(),
		"SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		r.log.Errorf("Failed to fetch user ID: %v", err)
		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		"SELECT item_type, quantity FROM inventory WHERE user_id = (SELECT id FROM users WHERE username=$1)", username)
	if err != nil {
		r.log.Errorf("Error fetching inventory: %v", err)
		return nil, err
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err = rows.Scan(&item.Type, &item.Quantity); err != nil {
			r.log.Errorf("Error scanning inventory row: %v", err)
			continue
		}
		inventory = append(inventory, item)
	}
	return inventory, nil
}

func (r *InventoryRepository) AddToInventory(username, itemType string, quantity int) error {
	// Получаем ID пользователя
	var userID int
	r.log.Infof("Attempting to add item to inventory: username=%s, item=%s", username, itemType)
	err := r.db.QueryRow(context.Background(),
		"SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		r.log.Errorf("User not found in inventory update: %s", username)
		return err
	}
	r.log.Infof("User found for inventory update: id=%d, username=%s", userID, username)

	// Добавляем предмет в инвентарь (или обновляем количество)
	r.log.Infof("Attempting to add item to inventory: username=%s, item=%s", username, itemType)
	_, err = r.db.Exec(context.Background(),
		`INSERT INTO inventory (user_id, item_type, quantity)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id, item_type)
         DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity`,
		userID, itemType, quantity)

	if err != nil {
		r.log.Errorf("Error adding user to inventory: %v", err)
		return fmt.Errorf("failed to add item to inventory: %w", err)
	}
	return nil
}
