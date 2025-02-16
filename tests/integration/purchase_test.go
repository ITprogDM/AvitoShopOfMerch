package integration

import (
	"ShopAvito/internal/handlers"
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

func fakeAuthMiddleware(c *gin.Context) {
	username := c.GetHeader("username")
	if username != "" {
		c.Set("username", username) // Устанавливаем username только если заголовок присутствует
	}
	c.Next()
}

func TestBuyItemAPI(t *testing.T) {
	// Инициализация тестовой базы данных
	db, err := testutils.InitTestDB()
	assert.NoError(t, err)
	defer db.Close()
	defer cleanupTestDB(db)

	// Создаем пользователя
	_, _, err = testutils.CreateTestUsers(db)
	assert.NoError(t, err)

	// Инициализация сервисов и обработчиков
	userRepo := repository.NewUserRepository(db, logrus.New())
	purchaseRepo := repository.NewPurchaseRepository(db, logrus.New())
	inventoryRepo := repository.NewInventoryRepository(db, logrus.New())
	inventoryService := services.NewInventoryService(inventoryRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, userRepo, inventoryRepo, logrus.New())
	purchaseHandler := handlers.NewPurchaseHandler(purchaseService, inventoryService, logrus.New())

	// Инициализация роутера
	router := gin.Default()
	router.Use(fakeAuthMiddleware)
	router.GET("/api/buy/:item", purchaseHandler.BuyItem)

	// Отправляем запрос на покупку товара
	req, _ := http.NewRequest("GET", "/api/buy/t-shirt", nil)
	req.Header.Set("username", "sender") // Устанавливаем покупателя в заголовке

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем ответ
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Purchase successful", response["message"])

	// Проверяем баланс пользователя
	balance, err := userRepo.GetUserBalance("sender")
	assert.NoError(t, err)
	assert.Equal(t, 920, balance) // 1000 - 80 = 920

	// Проверяем инвентарь пользователя
	inventory, err := inventoryRepo.GetInventory("sender")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(inventory))            // Должен быть один предмет
	assert.Equal(t, "t-shirt", inventory[0].Type) // Проверяем тип предмета
}
