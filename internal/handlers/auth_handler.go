package handlers

import (
	"ShopAvito/internal/models"
	"ShopAvito/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
	userService services.UserServiceInterface
	log         *logrus.Logger
}

func NewAuthHandler(authService services.AuthServiceInterface, userService services.UserServiceInterface, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		log:         log,
	}
}

func (h *AuthHandler) Authenticate(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверяем, существует ли пользователь
	exists, err := h.userService.UserExists(req.Username)
	if err != nil {
		h.log.Error("User verification error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var token string
	if exists {
		// Если пользователь существует, проверяем пароль и выдаем токен
		token, err = h.authService.Login(req.Username, req.Password)
		if err != nil {
			h.log.Error("Error log in:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
	} else {
		// Если пользователя нет, регистрируем его
		token, err = h.authService.Register(req.Username, req.Password)
		if err != nil {
			h.log.Error("Error register:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}
