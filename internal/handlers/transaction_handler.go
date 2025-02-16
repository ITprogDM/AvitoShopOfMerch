package handlers

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type TransactionHandler struct {
	transactionService services.TransactionServiceInterface
	log                *logrus.Logger
}

func NewTransactionHandler(transactionService services.TransactionServiceInterface, log *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		log:                log,
	}
}

func (h *TransactionHandler) SendCoins(c *gin.Context) {
	var req models.SendCoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("error occurred while binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	fromUser := c.MustGet("username").(string)
	err := h.transactionService.TransferCoins(fromUser, req.ToUser, req.Amount)
	if err != nil {
		h.log.Errorf("error occurred while sending coins: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
}
