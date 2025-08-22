package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-portfolio/rest-api/internal/auth"
	"github.com/go-portfolio/rest-api/internal/services"
)

// LoginHandler godoc
// @Summary      Авторизация пользователя
// @Description  Аутентификация пользователя и получение JWT токена
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  map[string]string  true  "Данные для входа"  example({"username": "user1", "password": "pass123"})
// @Success      200  {object}  map[string]string  "JWT токен"  example({"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."})
// @Failure      401  {string}  string  "unauthorized"
// @Router       /login [post]
func LoginHandler(userSvc services.UserService, jwtSecret string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

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
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}

