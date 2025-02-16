package testutils

import (
	"ShopAvito/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func InitTestDB() (*pgxpool.Pool, error) {
	connString := "postgres://user:password@localhost:5433/avito_shop_test?sslmode=disable"
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Применяем миграции
	if err = applyMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}

// applyMigrations применяет миграции к базе данных
func applyMigrations(db *pgxpool.Pool) error {
	// Пример миграции
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			balance INTEGER DEFAULT 1000
		);

		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_user INTEGER NOT NULL,
			to_user INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			timestamp TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (from_user) REFERENCES users(id),
			FOREIGN KEY (to_user) REFERENCES users(id)
		);

		CREATE TABLE IF NOT EXISTS purchases (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			item_name TEXT NOT NULL,
			price INTEGER NOT NULL,
			timestamp TIMESTAMP DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);

		CREATE TABLE IF NOT EXISTS inventory (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			item_type TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		    UNIQUE (user_id, item_type)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func CreateTestUsers(db *pgxpool.Pool) (senderID, receiverID int, err error) {
	// Хэшируем пароль для отправителя
	hashedSenderPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to hash sender password: %w", err)
	}

	// Создаем отправителя
	sender := models.User{
		Username: "sender",
		Password: string(hashedSenderPassword), // Сохраняем хэшированный пароль
		Balance:  1000,
	}
	err = db.QueryRow(context.Background(),
		"INSERT INTO users (username, password, balance) VALUES ($1, $2, $3) RETURNING id",
		sender.Username, sender.Password, sender.Balance,
	).Scan(&senderID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create sender: %w", err)
	}

	// Хэшируем пароль для получателя
	hashedReceiverPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to hash receiver password: %w", err)
	}

	// Создаем получателя
	receiver := models.User{
		Username: "receiver",
		Password: string(hashedReceiverPassword), // Сохраняем хэшированный пароль
		Balance:  500,
	}
	err = db.QueryRow(context.Background(),
		"INSERT INTO users (username, password, balance) VALUES ($1, $2, $3) RETURNING id",
		receiver.Username, receiver.Password, receiver.Balance,
	).Scan(&receiverID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create receiver: %w", err)
	}

	return senderID, receiverID, nil
}
