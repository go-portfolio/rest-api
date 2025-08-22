package models

import "time"

// User представляет пользователя
// swagger:model User
type User struct {
    // ID пользователя
    // example: 1
    ID int `json:"id"`

    // Логин пользователя
    // example: "user123"
    Username string `json:"username"`

    // Электронная почта пользователя
    // example: "user@example.com"
    Email string `json:"email"`

    // Хэш пароля пользователя
    // example: "$2a$10$E0NRl..."
    Password string `json:"password_hash"`

    // Дата создания пользователя в формате RFC3339
    // example: "2025-08-22T17:00:00Z"
    CreatedAt time.Time `json:"created_at"`
}

