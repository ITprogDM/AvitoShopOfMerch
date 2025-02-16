package integration

import (
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/models"
	"ShopAvito/internal/repository"
	"ShopAvito/internal/services"
	"ShopAvito/tests/testutils"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func cleanupTestDB(db *pgxpool.Pool) {
	_, err := db.Exec(context.Background(), `
		DROP TABLE IF EXISTS inventory;
		DROP TABLE IF EXISTS purchases;
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS users;
	`)
	if err != nil {
		log.Fatalf("Failed to clean up test database: %v", err)
	}
}

func TestSendCoinAPI(t *testing.T) {
	// Инициализация тестовой базы данных
	db, err := testutils.InitTestDB()
	assert.NoError(t, err)
	defer db.Close()
	defer cleanupTestDB(db)
	// Создаем двух пользователей
	_, _, err = testutils.CreateTestUsers(db)
	assert.NoError(t, err)

	// Инициализация сервисов и обработчиков
	userRepo := repository.NewUserRepository(db, logrus.New())
	transactionRepo := repository.NewTransactionRepository(db, logrus.New())
	transactionService := services.NewTransactionService(transactionRepo, userRepo, logrus.New())
	transactionHandler := handlers.NewTransactionHandler(transactionService, logrus.New())

	// Инициализация роутера
	router := gin.Default()
	router.Use(fakeAuthMiddleware)
	router.POST("/api/sendCoin", transactionHandler.SendCoins)

	// Создаем JSON-запрос для перевода монет
	requestBody := models.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	}
	jsonData, _ := json.Marshal(requestBody)

	// Отправляем запрос
	req, _ := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("username", "sender") // Устанавливаем отправителя в заголовке

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем ответ
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Transaction successful", response["message"])

	// Проверяем баланс отправителя и получателя
	senderBalance, err := userRepo.GetUserBalance("sender")
	assert.NoError(t, err)
	assert.Equal(t, 900, senderBalance) // 1000 - 100 = 900

	receiverBalance, err := userRepo.GetUserBalance("receiver")
	assert.NoError(t, err)
	assert.Equal(t, 600, receiverBalance) // 500 + 100 = 600
}
