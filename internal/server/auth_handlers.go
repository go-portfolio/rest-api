package server

import (
    "encoding/json"
    "net/http"
    "github.com/go-portfolio/rest-api/internal/auth"
    "github.com/go-portfolio/rest-api/internal/services"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func LoginHandler(userSvc services.UserService, jwtSecret string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid request", http.StatusBadRequest)
            return
        }

        user, err := userSvc.Authenticate(req.Username, req.Password)
        if err != nil {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }

        token, err := auth.GenerateToken(user.ID, jwtSecret)
        if err != nil {
            http.Error(w, "could not generate token", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}
