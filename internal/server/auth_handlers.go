package server

import (
    "encoding/json"
    "net/http"
    "github.com/go-portfolio/rest-api/internal/auth"
    "github.com/go-portfolio/rest-api/internal/services"
)


func LoginHandler(userSvc services.UserService, jwtSecret string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var creds struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        json.NewDecoder(r.Body).Decode(&creds)

        user, err := userSvc.Authenticate(creds.Username, creds.Password)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        token, _ := auth.GenerateToken(user.ID, jwtSecret)
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}

