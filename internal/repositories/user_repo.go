package repositories

import (
	"database/sql"
	"errors"

	"github.com/go-portfolio/rest-api/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

func FindUserByEmail(db *sql.DB, email string) (models.User, error) {
	var u models.User
	row := db.QueryRow(`SELECT id, email, password FROM users WHERE email = $1`, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}
	return u, nil
}

func CreateUser(db *sql.DB, email, hashed string) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id`, email, hashed).Scan(&id)
	return id, err
}
