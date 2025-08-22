package unit

import (
    "testing"

    "github.com/go-portfolio/rest-api/internal/models"
    "github.com/go-portfolio/rest-api/internal/services"
    "github.com/stretchr/testify/assert"
)

func TestMockUserService_Authenticate(t *testing.T) {
    // Создаём mock-сервис с тестовыми пользователями
    mockSvc := &services.MockUserService{
        Users: []models.User{
            {ID: 1, Username: "alex", Password: "password123"},
            {ID: 2, Username: "maria", Password: "secret456"},
        },
    }

    t.Run("correct credentials", func(t *testing.T) {
        user, err := mockSvc.Authenticate("alex", "password123")
        assert.NoError(t, err)
        assert.Equal(t, "alex", user.Username)
    })

    t.Run("wrong password", func(t *testing.T) {
        user, err := mockSvc.Authenticate("alex", "wrongpass")
        assert.Nil(t, user)
        assert.Error(t, err)
        assert.Equal(t, "invalid username or password", err.Error())
    })

    t.Run("nonexistent user", func(t *testing.T) {
        user, err := mockSvc.Authenticate("john", "password")
        assert.Nil(t, user)
        assert.Error(t, err)
        assert.Equal(t, "invalid username or password", err.Error())
    })
}
