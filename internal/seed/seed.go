package seed

import (
    "database/sql"
    "log"

    "golang.org/x/crypto/bcrypt"
)

// SeedUsers добавляет тестовых пользователей
func SeedUsers(db *sql.DB) {
    users := []struct {
        Username string
        Password string
        Email    string
    }{
        {"alex", "password123", "alex@example.com"},
        {"maria", "secret456", "maria@example.com"},
    }

    for _, u := range users {
        hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            log.Fatalf("failed to hash password: %v", err)
        }

        _, err = db.Exec(`
            INSERT INTO users (username, password_hash, email) 
            VALUES ($1, $2, $3)
            ON CONFLICT (username) DO NOTHING
        `, u.Username, string(hash), u.Email)
        if err != nil {
            log.Fatalf("failed to insert user %s: %v", u.Username, err)
        }
    }

    log.Println("Seed users done")
}
