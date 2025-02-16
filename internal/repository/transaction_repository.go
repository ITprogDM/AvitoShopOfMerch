package repository

import (
	"ShopAvito/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type TransactionRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewTransactionRepository(db *pgxpool.Pool, log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		db:  db,
		log: log,
	}
}

func (r *TransactionRepository) GetUserID(username string) (int, error) {
	var userID int
	err := r.db.QueryRow(context.Background(), "SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		r.log.Errorf("Failed to get user ID for username %s: %v", username, err)
		return 0, err
	}
	return userID, nil
}

// Перевод монет между пользователями
func (r *TransactionRepository) TransferCoins(fromUser, toUser string, amount int) error {
	fromUserID, err := r.GetUserID(fromUser)
	if err != nil {
		r.log.Errorf("Failed to get user ID for fromUser %s: %v", fromUser, err)
		return err
	}

	toUserID, err := r.GetUserID(toUser)
	if err != nil {
		r.log.Errorf("Failed to get user ID for toUser %s: %v", toUser, err)
		return err
	}
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		r.log.Error("Failed to begin transaction: ", err)
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

	// Вычитаем монеты у отправителя
	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, fromUserID)
	if err != nil {
		r.log.Error("Failed to update sender balance: ", err)
		return err
	}

	// Добавляем монеты получателю
	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toUserID)
	if err != nil {
		r.log.Error("Failed to update recipient balance: ", err)
		return err
	}

	// Записываем транзакцию в историю
	_, err = tx.Exec(context.Background(),
		"INSERT INTO transactions (from_user, to_user, amount) VALUES ($1, $2, $3)",
		fromUserID, toUserID, amount)
	if err != nil {
		r.log.Error("Failed to insert transaction record: ", err)
		return err
	}

	// Фиксируем транзакцию
	if err = tx.Commit(context.Background()); err != nil {
		r.log.Errorf("Failed to commit transaction: %v", err)
		return err
	}

	r.log.Info("Transaction committed successfully. Item added to inventory.")
	return nil
}

// Получение истории полученных монет
func (r *TransactionRepository) GetReceivedTransactions(username string) ([]models.Transaction, error) {
	userID, err := r.GetUserID(username)
	if err != nil {
		r.log.Errorf("Failed to get user ID for username %s: %v", username, err)
		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		`SELECT u1.username AS from_user, t.amount, t.timestamp 
        FROM transactions t
        JOIN users u1 ON t.from_user = u1.id
        WHERE t.to_user = $1`, userID)
	if err != nil {
		r.log.Error("Error fetching received transactions: ", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.FromUser, &transaction.Amount, &transaction.Time)
		if err != nil {
			r.log.Error("Error scanning received transaction: ", err)
			continue
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

// Получение истории отправленных монет
func (r *TransactionRepository) GetSentTransactions(username string) ([]models.Transaction, error) {
	userID, err := r.GetUserID(username)
	if err != nil {
		r.log.Errorf("Failed to get user ID for username %s: %v", username, err)
	}
	rows, err := r.db.Query(context.Background(),
		`SELECT u2.username AS to_user, t.amount, t.timestamp 
        FROM transactions t
        JOIN users u2 ON t.to_user = u2.id
        WHERE t.from_user = $1`, userID)
	if err != nil {
		r.log.Error("Error fetching sent transactions: ", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.ToUser, &transaction.Amount, &transaction.Time)
		if err != nil {
			r.log.Error("Error scanning sent transaction: ", err)
			continue
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
