package repository

import (
	"ShopAvito/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log *logrus.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (r *UserRepository) CreateUser(user models.User) error {
	_, err := r.db.Exec(context.Background(),
		"INSERT INTO users (username, password, balance) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Balance)
	return err
}

func (r *UserRepository) GetUserBalance(username string) (int, error) {
	var balance int
	err := r.db.QueryRow(context.Background(),
		"SELECT balance FROM users WHERE username = $1", username).Scan(&balance)
	return balance, err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(context.Background(),
		"SELECT id, username, password, balance FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Balance)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserBalance(username string, newBalance int) error {
	_, err := r.db.Exec(context.Background(),
		"UPDATE users SET balance = $1 WHERE username = $2", newBalance, username)
	return err
}

func (r *UserRepository) UserExists(username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		r.log.Errorf("The error is in the database request itself to check if the user exists: %s", err.Error())
	}
	return exists, err
}
