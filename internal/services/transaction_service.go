package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
	"errors"
	"github.com/sirupsen/logrus"
)

type TransactionService struct {
	transactionRepo repository.TransactionRepositoryInterface
	userRepo        repository.UserRepositoryInterface
	log             *logrus.Logger
}

func NewTransactionService(transactionRepo repository.TransactionRepositoryInterface, userRepo repository.UserRepositoryInterface, log *logrus.Logger) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		log:             log,
	}
}

// Перевод монет между пользователями
func (s *TransactionService) TransferCoins(fromUser, toUser string, amount int) error {
	balance, err := s.userRepo.GetUserBalance(fromUser)
	if err != nil {
		s.log.Errorf("Error getting user's balance: %v", err)
		return err
	}
	if balance < amount {
		return errors.New("insufficient funds")
	}

	if err = s.transactionRepo.TransferCoins(fromUser, toUser, amount); err != nil {
		s.log.Errorf("Error sending coins: %v", err)
		return err
	}
	return nil
}

// Получение истории полученных транзакций пользователя
func (s *TransactionService) GetReceivedTransactions(username string) ([]models.TransactionDetail, error) {
	received, err := s.transactionRepo.GetReceivedTransactions(username)
	if err != nil {
		s.log.Errorf("Error getting received transctions: %v", err)
		return nil, err
	}

	var details []models.TransactionDetail
	for _, t := range received {
		details = append(details, models.TransactionDetail{
			FromUser: t.FromUser,
			Amount:   t.Amount,
		})
	}
	return details, nil
}

// Получение истории отправленных транзакций пользователя
func (s *TransactionService) GetSentTransactions(username string) ([]models.TransactionDetail, error) {
	sent, err := s.transactionRepo.GetSentTransactions(username)
	if err != nil {
		s.log.Errorf("Error getting sent transactions: %v", err)
		return nil, err
	}

	var details []models.TransactionDetail
	for _, t := range sent {
		details = append(details, models.TransactionDetail{
			ToUser: t.ToUser,
			Amount: t.Amount,
		})
	}
	return details, nil
}
