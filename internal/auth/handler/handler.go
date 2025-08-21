package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-portfolio/rest-api/internal/auth/jwt"  // Импортируем пакет JWT
)

// Обработчик для логина
func LoginHandler(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверка логина и пароля (можно заменить на реальную логику проверки)
	if loginRequest.Username == "user" && loginRequest.Password == "password" {
		token, err := jwt.GenerateJWT(loginRequest.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
}

// Защищённый маршрут
func ProtectedHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// Валидация JWT
	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Protected data", "user": claims.Username})
}
