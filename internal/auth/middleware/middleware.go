package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-portfolio/rest-api/internal/auth/jwt"  // Импортируем пакет JWT
)

// Middleware для защиты маршрутов с токенами
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}

		// Валидация JWT
		claims, err := jwt.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Добавление информации о пользователе в контекст
		c.Set("user", claims.Username)
		c.Next()
	}
}
