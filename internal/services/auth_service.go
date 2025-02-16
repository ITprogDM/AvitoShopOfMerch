package services

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userRepo  repository.UserRepositoryInterface
	secretKey string
	log       *logrus.Logger
}

func NewAuthService(userRepo repository.UserRepositoryInterface, secretKey string, log *logrus.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		secretKey: secretKey,
		log:       log,
	}
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	rToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.log.Errorf("Error while signing token: %v", err)
		return "", err
	}
	return rToken, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		s.log.Errorf("Token validation error: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		s.log.Error("Invalid token or claims")
		return nil, errors.New("invalid token")
	}

	s.log.Infof("Token validated successfully. Username: %s", claims.Username)
	return claims, nil
}

// Логин (проверка пароля и выдача токена)
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		s.log.Errorf("Error login while getting user by username: %v", err)
		return "", errors.New("user not found")
	}

	// Проверяем пароль
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.log.Errorf("Error login while comparing password: %v", err)
		return "", errors.New("invalid password")
	}

	// Генерируем токен
	return s.GenerateToken(user.Username)
}

// Регистрация (создание пользователя и выдача токена)
func (s *AuthService) Register(username, password string) (string, error) {
	// Проверяем, есть ли такой пользователь
	exists, err := s.userRepo.UserExists(username)
	if err != nil {
		s.log.Errorf("Error finding user: %v", err)
		return "", err
	}
	if exists {
		return "", errors.New("user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Errorf("Error hashing password: %v", err)
		return "", err
	}

	// Создаём пользователя
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Balance:  1000,
	}
	err = s.userRepo.CreateUser(user)
	if err != nil {
		s.log.Errorf("Error creating user: %v", err)
		return "", err
	}

	// Генерируем токен
	return s.GenerateToken(username)
}
