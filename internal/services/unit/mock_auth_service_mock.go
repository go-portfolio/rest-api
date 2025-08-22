package unit


import (
    "testing"

    _ "github.com/lib/pq" // Подключаем драйвер PostgreSQL
    "github.com/stretchr/testify/assert"
    "github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services" 
)


// ------------------------
// Unit-тест для MockUserService (аутентификация пользователей)
// ------------------------
func TestUnitMockUserService_Authenticate(t *testing.T) {
    // Создаём mock-сервис с тестовыми пользователями
    mockSvc := &services.MockUserService{
        Users: []models.User{
            {ID: 1, Username: "alex", Password: "password123"},
            {ID: 2, Username: "maria", Password: "secret456"},
        },
    }

    // Тест успешной аутентификации с правильными данными
    t.Run("correct credentials", func(t *testing.T) {
        user, err := mockSvc.Authenticate("alex", "password123")
        assert.NoError(t, err)              // Ошибок быть не должно
        assert.Equal(t, "alex", user.Username) // Проверяем имя пользователя
    })

    // Тест с неправильным паролем
    t.Run("wrong password", func(t *testing.T) {
        user, err := mockSvc.Authenticate("alex", "wrongpass")
        assert.Nil(t, user)                 // Пользователь не должен вернуться
        assert.Error(t, err)                // Должна быть ошибка
        assert.Equal(t, "invalid username or password", err.Error())
    })

    // Тест с несуществующим пользователем
    t.Run("nonexistent user", func(t *testing.T) {
        user, err := mockSvc.Authenticate("john", "password")
        assert.Nil(t, user)                 // Пользователь не найден
        assert.Error(t, err)                // Должна быть ошибка
        assert.Equal(t, "invalid username or password", err.Error())
    })
}

