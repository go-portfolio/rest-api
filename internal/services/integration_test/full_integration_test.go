package integration_test

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-portfolio/rest-api/internal/config"
    "github.com/go-portfolio/rest-api/internal/services"
    "github.com/go-portfolio/rest-api/internal/models"
    "github.com/go-portfolio/rest-api/internal/server"
    "github.com/go-portfolio/rest-api/internal/auth"
    "github.com/stretchr/testify/assert"
    _ "github.com/lib/pq" // PostgreSQL драйвер
)

func TestFullIntegration(t *testing.T) {
    // --------------------------
    // Загружаем конфигурацию проекта
    // cfg содержит DSN для подключения к БД, JWT секреты и др.
    // --------------------------
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err)
    }

    // --------------------------
    // Подключаемся к тестовой базе PostgreSQL
    // --------------------------
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close() // Закрываем соединение после теста

    // --------------------------
    // Создаём реальные сервисы для работы с задачами
    // --------------------------
    taskSvc := services.NewPostgresTaskService(db)

    // --------------------------
    // Создаём тестового пользователя в базе
    // Можно хранить пароль в явном виде для упрощения теста
    // --------------------------
    _, err = db.Exec(`INSERT INTO users (username, password_hash) VALUES ($1, $2)
                      ON CONFLICT (username) DO NOTHING`,
        "testuser", "password123") // Для реального проекта лучше использовать bcrypt
    if err != nil {
        t.Fatal(err)
    }

    // --------------------------
    // 1. Генерация JWT для пользователя
    // --------------------------
    token, err := auth.GenerateToken(1, cfg.Jwt.JwtSecretKey) // user_id = 1
    if err != nil {
        t.Fatal(err)
    }

    // --------------------------
    // 2. Тестируем HTTP POST /tasks с JWT
    // Создаём тестовый сервер с middleware проверки токена
    // --------------------------
    ts := httptest.NewServer(
        auth.VerifyToken(cfg.Jwt.JwtSecretKey)(
            http.HandlerFunc(server.TasksHandler(taskSvc)),
        ),
    )
    defer ts.Close()

    // Создаём новую задачу
    newTask := models.Task{Title: "Integration Task", Status: "todo"}
    body, _ := json.Marshal(newTask) // Преобразуем в JSON

    // Создаём POST-запрос с JWT в заголовке Authorization
    req, _ := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    // Отправляем запрос
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close() // Закрываем тело ответа

    // Проверяем, что сервер вернул HTTP 200 OK
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Декодируем ответ в структуру задачи
    var createdTask models.Task
    json.NewDecoder(resp.Body).Decode(&createdTask)
    assert.Equal(t, newTask.Title, createdTask.Title)
    assert.Equal(t, newTask.Status, createdTask.Status)

    // --------------------------
    // 3. Проверка, что задача действительно сохранилась в базе
    // --------------------------
    row := db.QueryRow(`SELECT id, title, status FROM tasks WHERE id=$1`, createdTask.ID)
    var taskFromDB models.Task
    err = row.Scan(&taskFromDB.ID, &taskFromDB.Title, &taskFromDB.Status)
    if err != nil {
        t.Fatal(err)
    }

    // Сравниваем ID созданной задачи и задачи в базе
    assert.Equal(t, createdTask.ID, taskFromDB.ID)
}
