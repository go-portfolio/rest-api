package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password_hash"`
	CreatedAt time.Time `db:"created_at"` 
}
