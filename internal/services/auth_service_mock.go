package services

import (
	"errors"

	"github.com/go-portfolio/rest-api/internal/models"
)

type MockUserService struct {
    Users []models.User // можно хранить пароль в явном виде для простоты
}

func (m *MockUserService) Authenticate(username, password string) (*models.User, error) {
    for _, u := range m.Users {
        if u.Username == username && u.Password == password { // в mock можно хранить plain password
            return &u, nil
        }
    }
    return nil, errors.New("invalid username or password")
}
