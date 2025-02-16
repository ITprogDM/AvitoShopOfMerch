package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

type StubTransactionRepository struct {
	GetUserIDFunc               func(username string) (int, error)
	TransferCoinsFunc           func(fromUser, toUser string, amount int) error
	GetReceivedTransactionsFunc func(username string) ([]models.Transaction, error)
	GetSentTransactionsFunc     func(username string) ([]models.Transaction, error)
}

func (r *StubTransactionRepository) GetUserID(username string) (int, error) {
	return r.GetUserIDFunc(username)
}

func (s *StubTransactionRepository) TransferCoins(fromUser, toUser string, amount int) error {
	return s.TransferCoinsFunc(fromUser, toUser, amount)
}

func (s *StubTransactionRepository) GetReceivedTransactions(username string) ([]models.Transaction, error) {
	return s.GetReceivedTransactionsFunc(username)
}

func (s *StubTransactionRepository) GetSentTransactions(username string) ([]models.Transaction, error) {
	return s.GetSentTransactionsFunc(username)
}

// StubUserRepository - заглушка для UserRepository
type StubUserRepositoryForTransactions struct {
	GetUserBalanceFunc func(username string) (int, error)
}

func (s *StubUserRepositoryForTransactions) GetUserBalance(username string) (int, error) {
	return s.GetUserBalanceFunc(username)
}

func TestTransactionService_TransferCoins(t *testing.T) {
	// Создаем заглушки для репозиториев
	stubTransactionRepo := &StubTransactionRepository{
		TransferCoinsFunc: func(fromUser, toUser string, amount int) error {
			if fromUser == "sender" && toUser == "receiver" && amount == 100 {
				return nil
			}
			return errors.New("failed to transfer coins")
		},
	}

	stubUserRepo := &StubUserRepository{
		GetUserBalanceFunc: func(username string) (int, error) {
			if username == "sender" {
				return 1000, nil // У отправителя достаточно баланса
			}
			return 0, errors.New("user not found")
		},
	}

	// Инициализируем логгер с отключенным выводом
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Отключаем вывод логов в консоль

	// Создаем TransactionService с заглушками и логгером
	transactionService := services.NewTransactionService(stubTransactionRepo, stubUserRepo, logger)

	// Тест на успешный перевод монет
	err := transactionService.TransferCoins("sender", "receiver", 100)
	assert.NoError(t, err)

	// Тест на недостаточный баланс
	stubUserRepo.GetUserBalanceFunc = func(username string) (int, error) {
		return 50, nil // У отправителя недостаточно баланса
	}
	err = transactionService.TransferCoins("sender", "receiver", 100)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	// Тест на ошибку при переводе монет
	stubUserRepo.GetUserBalanceFunc = func(username string) (int, error) {
		return 1000, nil // У отправителя достаточно баланса
	}
	stubTransactionRepo.TransferCoinsFunc = func(fromUser, toUser string, amount int) error {
		return errors.New("failed to transfer coins") // Симулируем ошибку при переводе
	}
	err = transactionService.TransferCoins("sender", "receiver", 100)
	assert.Error(t, err)
	assert.Equal(t, "failed to transfer coins", err.Error())
}

func TestTransactionService_GetReceivedTransactions(t *testing.T) {
	// Создаем заглушку для TransactionRepository
	stubTransactionRepo := &StubTransactionRepository{
		GetReceivedTransactionsFunc: func(username string) ([]models.Transaction, error) {
			if username == "receiver" {
				return []models.Transaction{
					{
						FromUser: "sender",
						Amount:   100,
						Time:     time.Now(),
					},
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	// Инициализируем логгер с отключенным выводом
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Отключаем вывод логов в консоль

	// Создаем TransactionService с заглушкой и логгером
	transactionService := services.NewTransactionService(stubTransactionRepo, nil, logger)

	// Тест на успешное получение списка полученных транзакций
	transactions, err := transactionService.GetReceivedTransactions("receiver")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, "sender", transactions[0].FromUser)
	assert.Equal(t, 100, transactions[0].Amount)

	// Тест на несуществующего пользователя
	_, err = transactionService.GetReceivedTransactions("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestTransactionService_GetSentTransactions(t *testing.T) {
	// Создаем заглушку для TransactionRepository
	stubTransactionRepo := &StubTransactionRepository{
		GetSentTransactionsFunc: func(username string) ([]models.Transaction, error) {
			if username == "sender" {
				return []models.Transaction{
					{
						ToUser: "receiver",
						Amount: 100,
						Time:   time.Now(),
					},
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	// Инициализируем логгер с отключенным выводом
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Отключаем вывод логов в консоль

	// Создаем TransactionService с заглушкой и логгером
	transactionService := services.NewTransactionService(stubTransactionRepo, nil, logger)

	// Тест на успешное получение списка отправленных транзакций
	transactions, err := transactionService.GetSentTransactions("sender")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, "receiver", transactions[0].ToUser)
	assert.Equal(t, 100, transactions[0].Amount)

	// Тест на несуществующего пользователя
	_, err = transactionService.GetSentTransactions("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}
