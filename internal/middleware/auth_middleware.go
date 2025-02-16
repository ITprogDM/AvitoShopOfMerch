package middleware

import (
	"ShopAvito/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func AuthMiddleware(authService *services.AuthService, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		log.Info("Authorization Header:", tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Убираем "Bearer " перед токеном
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString = parts[1] // Оставляем только сам токен

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			log.Info("Token validation error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Infof("Token valid. Username extracted: %s", claims.Username)
		// Передаем username в контекст запроса
		c.Set("username", claims.Username)
		c.Next()
	}
}
