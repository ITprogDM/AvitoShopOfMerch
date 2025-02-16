package handlers

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	userService        services.UserServiceInterface
	transactionService services.TransactionServiceInterface
	inventoryService   services.InventoryServiceInterface
	log                *logrus.Logger
}

func NewUserHandler(userService services.UserServiceInterface, transactionService services.TransactionServiceInterface, inventoryService services.InventoryServiceInterface, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService:        userService,
		transactionService: transactionService,
		inventoryService:   inventoryService,
		log:                log,
	}
}

func (h *UserHandler) GetBalance(c *gin.Context) {
	username := c.MustGet("username").(string)
	balance, err := h.userService.GetBalance(username)
	if err != nil {
		h.log.Errorf("Error getting balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {

	usernameAny, exists := c.Get("username")
	if !exists {
		// Если username отсутствует, возвращаем статус 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Приводим username к типу string
	username, ok := usernameAny.(string)
	if !ok {
		// Если username не является строкой, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username type"})
		return
	}

	h.log.Infof("Fetching info for user: %s", username)

	balance, err := h.userService.GetBalance(username)
	if err != nil {
		h.log.Errorf("Error fetching balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		return
	}

	h.log.Infof("Fetching transactions and inventory for user: %s", username)
	received, err := h.transactionService.GetReceivedTransactions(username)
	if err != nil {
		h.log.Errorf("Error fetching received transactions: %v", err)
		received = []models.TransactionDetail{} // <-- Теперь не null!
	}

	h.log.Infof("Fetching transactions and inventory for user: %s", username)
	sent, err := h.transactionService.GetSentTransactions(username)
	if err != nil {
		h.log.Errorf("Error fetching sent transactions: %v", err)
		sent = []models.TransactionDetail{} // <-- Теперь не null!
	}

	h.log.Infof("Fetching transactions and inventory for user: %s", username)
	inventory, err := h.inventoryService.GetInventory(username)
	if err != nil {
		h.log.Errorf("Error fetching inventory: %v", err)
		inventory = []models.InventoryItem{} // <-- Теперь не null!
	}

	c.JSON(http.StatusOK, gin.H{
		"coins":     balance,
		"inventory": inventory,
		"coinHistory": gin.H{
			"received": received,
			"sent":     sent,
		},
	})
}
