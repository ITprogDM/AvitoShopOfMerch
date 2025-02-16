package handlers

import (
	"ShopAvito/internal/middleware"
	"ShopAvito/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterRoutes(userService *services.UserService, transactionService *services.TransactionService, purchaseService *services.PurchaseService, authService *services.AuthService, invenService *services.InventoryService, log *logrus.Logger) *gin.Engine {
	authHandler := NewAuthHandler(authService, userService, log)
	userHandler := NewUserHandler(userService, transactionService, invenService, log)
	transactionHandler := NewTransactionHandler(transactionService, log)
	purchaseHandler := NewPurchaseHandler(purchaseService, invenService, log)

	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/auth", authHandler.Authenticate)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService, log))
		{
			protected.GET("/info", userHandler.GetUserInfo)
			protected.POST("/sendCoin", transactionHandler.SendCoins)
			protected.GET("/buy/:item", purchaseHandler.BuyItem)
		}
	}

	return router
}
