package services

import (
	"database/sql"
	"errors"

	"github.com/go-portfolio/rest-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Интерфейс для работы с пользователями
type UserService interface {
    Authenticate(username, password string) (*models.User, error)
	FindUserByEmail(db *sql.DB, email string) (models.User, error)
	CreateUser(db *sql.DB, email, hashed string) (int, error)	
}

// Реализация UserService для Postgres
type PostgresUserService struct {
    DB *sql.DB
	Users []models.User
}

// Конструктор
func NewPostgresUserService(db *sql.DB) *PostgresUserService {
    return &PostgresUserService{DB: db}
}



func (p *PostgresUserService) Authenticate(username, password string) (*models.User, error) {
    var user models.User
    err := p.DB.QueryRow(`SELECT id, username, password_hash FROM users WHERE username=$1`, username).
        Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return nil, err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid password")
    }

    return &user, nil
}


var ErrUserNotFound = errors.New("user not found")

func (p *PostgresUserService) FindUserByEmail(db *sql.DB, email string) (models.User, error) {
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

func (p *PostgresUserService) CreateUser(db *sql.DB, email, hashed string) (int, error) {
	var id int
	err := db.QueryRow(`INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id`, email, hashed).Scan(&id)
	return id, err
}
