package services

import (
	"ShopAvito/internal/repository"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
	log      *logrus.Logger
}

func NewUserService(userRepo repository.UserRepositoryInterface, log *logrus.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		log:      log,
	}
}

func (s *UserService) UserExists(username string) (bool, error) {
	exists, err := s.userRepo.UserExists(username)
	if err != nil {
		s.log.Infof("Error checking if user exists: %s", err.Error())
		return false, err
	}
	return exists, nil
}

// Получение баланса пользователя
func (s *UserService) GetBalance(username string) (int, error) {
	return s.userRepo.GetUserBalance(username)
}
