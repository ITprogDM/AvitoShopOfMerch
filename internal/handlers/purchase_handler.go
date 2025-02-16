package handlers

import (
	"ShopAvito/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PurchaseHandler struct {
	purchaseService  services.PurchaseServiceInterface
	inventoryService services.InventoryServiceInterface
	log              *logrus.Logger
}

func NewPurchaseHandler(purchaseService services.PurchaseServiceInterface, inventoryService services.InventoryServiceInterface, log *logrus.Logger) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService:  purchaseService,
		inventoryService: inventoryService,
		log:              log,
	}
}

func (h *PurchaseHandler) BuyItem(c *gin.Context) {
	username := c.MustGet("username").(string)
	item := c.Param("item")

	// Цены товаров
	prices := map[string]int{
		"t-shirt": 80, "cup": 20, "book": 50,
		"pen": 10, "powerbank": 200, "hoody": 300,
		"umbrella": 200, "socks": 10, "wallet": 50, "pink-hoody": 500,
	}

	price, exists := prices[item]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item"})
		return
	}

	err := h.purchaseService.BuyItem(username, item, price)
	if err != nil {
		h.log.Errorf("Error buying item: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase successful"})
}
