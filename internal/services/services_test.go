package services

import (
    "database/sql"
    "testing"
    "log"

    "github.com/go-portfolio/rest-api/internal/config" 
    _ "github.com/lib/pq" // Подключаем драйвер PostgreSQL
    "github.com/stretchr/testify/assert"
    "github.com/go-portfolio/rest-api/internal/models"
)

// ------------------------
// Интеграционный тест: получение задач из БД
// ------------------------
func TestGetTasks(t *testing.T) {
    // Загружаем конфиг из корня проекта (DSN, настройки БД и т.п.)
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err) // Завершаем тест при ошибке загрузки конфига
    }

    // Подключаемся к базе данных PostgreSQL
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close() // Закрываем соединение после теста

    // Создаем сервис для работы с задачами
    svc := NewPostgresTaskService(db)
    
    // Получаем список задач
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatal(err)
    }

    // Проверяем, что список задач не пустой
    if len(tasks) == 0 {
        t.Error("expected at least one task")
    }
}

// ------------------------
// Интеграционный тест: создание и проверка задачи
// ------------------------
func TestCreateAndGetTasks(t *testing.T) {
    // Загружаем конфиг проекта
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err)
    }

    // Подключаемся к БД
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Создаем сервис для работы с задачами
    svc := NewPostgresTaskService(db)

    // Создаем тестовую задачу
    id, err := svc.CreateTask("Test task", "todo")
    if err != nil {
        t.Fatalf("Failed to create task: %v", err)
    }

    // Получаем список задач после создания
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatalf("Failed to get tasks: %v", err)
    }

    // Проверяем, что созданная задача присутствует в списке
    found := false
    for _, task := range tasks {
        if task.ID == id {
            found = true
            break
        }
    }

    if !found {
        t.Errorf("Таск не найден после создания")
    }
}

// ------------------------
// Unit-тест для MockUserService (аутентификация пользователей)
// ------------------------
func TestMockUserService_Authenticate(t *testing.T) {
    // Создаём mock-сервис с тестовыми пользователями
    mockSvc := &MockUserService{
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
