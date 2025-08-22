// models/user.go
package models

import "time"

// User представляет пользователя системы
// swagger:model User
type User struct {
    // ID пользователя
    // example: 1
    ID int `json:"id"`

    // Логин пользователя
    // example: "user123"
    // min length: 3
    // max length: 20
    Username string `json:"username" validate:"required,min=3,max=20"`

    // Электронная почта пользователя
    // example: "user@example.com"
    // format: email
    Email string `json:"email" validate:"required,email"`

    // Хэш пароля пользователя
    // example: "$2a$10$E0NRl..."
    // min length: 6
    Password string `json:"password_hash" validate:"required,min=6"`

    // Дата создания пользователя в формате RFC3339
    // example: "2025-08-22T17:00:00Z"
    CreatedAt time.Time `json:"created_at"`
}
