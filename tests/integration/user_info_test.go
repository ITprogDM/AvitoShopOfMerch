package integration

import (
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
	"ShopAvito/internal/services"
	"ShopAvito/tests/testutils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserInfoAPI(t *testing.T) {
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
	transactionRepo := repository.NewTransactionRepository(db, logrus.New())
	inventoryRepo := repository.NewInventoryRepository(db, logrus.New())
	transactionService := services.NewTransactionService(transactionRepo, userRepo, logrus.New())
	inventoryService := services.NewInventoryService(inventoryRepo)
	userService := services.NewUserService(userRepo, logrus.New())
	userHandler := handlers.NewUserHandler(userService, transactionService, inventoryService, logrus.New())

	// Инициализация роутера
	router := gin.Default()
	router.Use(fakeAuthMiddleware) // Используем фейковую аутентификацию
	router.GET("/api/info", userHandler.GetUserInfo)

	// Тест 1: Успешное получение информации о пользователе
	t.Run("Successful get user info", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/info", nil)
		req.Header.Set("username", "sender") // Устанавливаем username вручную

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.InfoResponse
		if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Equal(t, 1000, response.Coins)                  // Начальный баланс
		assert.Equal(t, 0, len(response.Inventory))            // Инвентарь пуст
		assert.Equal(t, 0, len(response.CoinHistory.Received)) // Нет полученных транзакций
		assert.Equal(t, 0, len(response.CoinHistory.Sent))     // Нет отправленных транзакций
	})

	// Тест 2: Пользователь не авторизован (без заголовка username)
	t.Run("Unauthorized user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/info", nil)
		// Не устанавливаем заголовок username

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
