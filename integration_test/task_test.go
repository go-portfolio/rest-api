package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "database/sql"

    "github.com/go-portfolio/rest-api/internal/services"
    "github.com/go-portfolio/rest-api/internal/server"
    "github.com/go-portfolio/rest-api/internal/config" 
    "github.com/stretchr/testify/assert"
    _ "github.com/lib/pq" // Для работы с PostgreSQL драйвером
    "github.com/go-portfolio/rest-api/internal/models"
)


// Интеграционный тест для создания новой задачи через HTTP
func TestCreateAndGetTask(t *testing.T) {
    // Загружаем конфигурацию из файла
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err) // Завершаем тест с ошибкой, если конфиг не загрузился
    }

    // Подключаемся к базе данных PostgreSQL
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        t.Fatal(err) // Завершаем тест с ошибкой, если не удалось подключиться к БД
    }
    defer db.Close() // Закрываем соединение с БД по завершению теста

    // Создаем реальный сервис для работы с задачами через PostgreSQL
    realSvc := services.NewPostgresTaskService(db)

    // Поднимаем тестовый HTTP сервер
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        handler := server.TasksHandler(realSvc) // Создаем обработчик задач с реальным сервисом
        handler.ServeHTTP(w, r) // Обрабатываем HTTP-запрос
    }))
    defer ts.Close() // Закрываем сервер после завершения теста

    // Данные для новой задачи
    newTask := models.Task{
        Title:  "New Task", // Заголовок задачи
        Status: "pending",  // Статус задачи
    }

    // Преобразуем данные задачи в JSON
    taskJSON, err := json.Marshal(newTask)
    if err != nil {
        t.Fatal(err) // Завершаем тест с ошибкой, если не удалось преобразовать в JSON
    }

    // Отправляем POST-запрос на создание новой задачи
    res, err := http.Post(ts.URL+"/tasks", "application/json", bytes.NewBuffer(taskJSON))
    if err != nil {
        t.Fatal(err) // Завершаем тест с ошибкой, если запрос не удался
    }
    defer res.Body.Close() // Закрываем тело ответа по завершению теста

    // Проверяем, что статус код ответа равен 200 (OK)
    assert.Equal(t, http.StatusOK, res.StatusCode)

    // Декодируем ответ в структуру задачи
    var createdTask models.Task
    err = json.NewDecoder(res.Body).Decode(&createdTask)
    if err != nil {
        t.Fatal(err) // Завершаем тест с ошибкой, если не удалось декодировать ответ
    }

    // Проверяем, что данные в ответе совпадают с отправленными
    assert.Equal(t, newTask.Title, createdTask.Title)  // Проверка заголовка задачи
    assert.Equal(t, newTask.Status, createdTask.Status) // Проверка статуса задачи

    // Проверяем, что задача получила ID (оно должно быть больше нуля)
    assert.Greater(t, createdTask.ID, 0)

    // Очищаем базу данных от тестовых данных после выполнения теста
    defer func() {
        _, err := db.Exec("DELETE FROM tasks WHERE title = 'New Task'")  // Очистка данных
        if err != nil {
            t.Fatal("Error cleaning up test data: ", err)
        }
    }()
}
