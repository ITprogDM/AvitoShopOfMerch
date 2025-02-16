package integration

import (
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
	"ShopAvito/internal/services"
	"ShopAvito/tests/testutils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthAPI(t *testing.T) {
	// Инициализация тестовой базы данных
	db, err := testutils.InitTestDB()
	assert.NoError(t, err)
	defer db.Close()
	defer cleanupTestDB(db)

	// Создаем тестового пользователя
	_, _, err = testutils.CreateTestUsers(db)
	assert.NoError(t, err)

	// Инициализация сервисов и обработчиков
	userRepo := repository.NewUserRepository(db, logrus.New())
	authService := services.NewAuthService(userRepo, "secret-key", logrus.New())
	userService := services.NewUserService(userRepo, logrus.New())
	authHandler := handlers.NewAuthHandler(authService, userService, logrus.New())

	// Инициализация роутера
	router := gin.Default()
	router.POST("/api/auth", authHandler.Authenticate)

	// Тест 1: Успешная аутентификация
	t.Run("Successful authentication", func(t *testing.T) {
		authRequest := models.AuthRequest{
			Username: "sender",
			Password: "password",
		}
		jsonData, _ := json.Marshal(authRequest)

		req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.AuthResponse
		if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		assert.NotEmpty(t, response.Token)
	})

	// Тест 2: Неверный пароль
	t.Run("Invalid password", func(t *testing.T) {
		authRequest := models.AuthRequest{
			Username: "sender",
			Password: "wrong-password",
		}
		jsonData, _ := json.Marshal(authRequest)

		req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Тест 3: Регистрация нового пользователя
	t.Run("Register new user", func(t *testing.T) {
		authRequest := models.AuthRequest{
			Username: "new-user",
			Password: "new-password",
		}
		jsonData, _ := json.Marshal(authRequest)

		req, _ := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.AuthResponse
		if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		assert.NotEmpty(t, response.Token)
	})
}
