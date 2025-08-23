package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-portfolio/rest-api/internal/auth"
	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services"
)

// LoginHandler godoc
// @Summary      Авторизация пользователя
// @Description  Аутентификация пользователя и получение JWT токена
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  LoginRequest  true  "Данные для входа"
// @Success      200  {object}  LoginResponse  "JWT токен и данные пользователя"
// @Failure      400  {object}  map[string]string  "Некорректный JSON"
// @Failure      401  {object}  map[string]string  "Неверные учетные данные"
// @Router       /login [post]
func LoginHandler(userSvc services.UserService, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		user, err := userSvc.Authenticate(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token, _ := auth.GenerateToken(user.ID, jwtSecret)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: token,
			User:  *user,
		})
	}
}

// LoginRequest модель запроса для Swagger
// swagger:model LoginRequest
type LoginRequest struct {
	// Логин пользователя
	// example: user1
	Username string `json:"username" validate:"required"`
	// Пароль пользователя
	// example: pass123
	Password string `json:"password" validate:"required"`
}

// LoginResponse модель ответа для Swagger
// swagger:model LoginResponse
type LoginResponse struct {
	// JWT токен
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
	// Данные пользователя
	User models.User `json:"user"`
}
