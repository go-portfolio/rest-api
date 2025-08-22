package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
    "strconv"

	"github.com/go-portfolio/rest-api/internal/auth"
	"github.com/go-portfolio/rest-api/internal/config"
	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/server"
	"github.com/go-portfolio/rest-api/internal/services"
	_ "github.com/lib/pq" // Для работы с PostgreSQL драйвером
	"github.com/stretchr/testify/assert"
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

	userID := 1
	// Генерируем тестовый токен для пользователя с ID=1
	token, _ := auth.GenerateToken(userID, cfg.Jwt.JwtSecretKey)
	ts := httptest.NewServer(auth.VerifyToken(cfg.Jwt.JwtSecretKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.TasksHandler(realSvc).ServeHTTP(w, r)
	})))
	defer ts.Close() // Закрываем сервер после завершения теста

	// Данные для новой задачи
	newTask := models.Task{
		Title:  "New Task", // Заголовок задачи
		Status: "pending",  // Статус задачи
		UserID: userID,
	}

	// Преобразуем данные задачи в JSON
	taskJSON, err := json.Marshal(newTask)
	if err != nil {
		t.Fatal(err) // Завершаем тест с ошибкой, если не удалось преобразовать в JSON
	}

	req, _ := http.NewRequest("POST", ts.URL+"/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)

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
	assert.Equal(t, newTask.Title, createdTask.Title)   // Проверка заголовка задачи
	assert.Equal(t, newTask.Status, createdTask.Status) // Проверка статуса задачи

	// Проверяем, что задача получила ID (оно должно быть больше нуля)
	assert.Greater(t, createdTask.ID, 0)

	// Очищаем базу данных от тестовых данных после выполнения теста
	defer func() {
		_, err := db.Exec("DELETE FROM tasks WHERE title = 'New Task'") // Очистка данных
		if err != nil {
			t.Fatal("Error cleaning up test data: ", err)
		}
	}()
}


// Интеграционный тест для обновления задачи через HTTP
func TestUpdateTask(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем сервис для работы с задачами
	realSvc := services.NewPostgresTaskService(db)

	userID := 1
	// Генерируем тестовый токен
	token, _ := auth.GenerateToken(userID, cfg.Jwt.JwtSecretKey)

	// Создаём HTTP тестовый сервер с авторизацией
	ts := httptest.NewServer(auth.VerifyToken(cfg.Jwt.JwtSecretKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.TasksHandler(realSvc).ServeHTTP(w, r)
	})))
	defer ts.Close()

	// --- Сначала создаём новую задачу ---
	newTask := models.Task{
		Title:  "Task to Update",
		Status: "pending",
		UserID: userID,
	}
	taskJSON, _ := json.Marshal(newTask)

	reqCreate, _ := http.NewRequest("POST", ts.URL+"/tasks", bytes.NewBuffer(taskJSON))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("Authorization", "Bearer "+token)
	resCreate, err := http.DefaultClient.Do(reqCreate)
	if err != nil {
		t.Fatal(err)
	}
	defer resCreate.Body.Close()

	var createdTask models.Task
	json.NewDecoder(resCreate.Body).Decode(&createdTask)

	// --- Теперь обновляем задачу ---
	updatedTask := models.Task{
		Title:  "Updated Task",
		Status: "done",
		UserID: userID,
	}
	updateJSON, _ := json.Marshal(updatedTask)

	url := ts.URL+"/tasks/"+strconv.Itoa(createdTask.ID)
	reqUpdate, _ := http.NewRequest("PUT", url, bytes.NewBuffer(updateJSON))
	reqUpdate.Header.Set("Content-Type", "application/json")
	reqUpdate.Header.Set("Authorization", "Bearer "+token)
	resUpdate, err := http.DefaultClient.Do(reqUpdate)
	if err != nil {
		t.Fatal(err)
	}
	defer resUpdate.Body.Close()

	assert.Equal(t, http.StatusOK, resUpdate.StatusCode)

	var returnedTask models.Task
	json.NewDecoder(resUpdate.Body).Decode(&returnedTask)

	// Проверяем, что данные обновились
	assert.Equal(t, updatedTask.Title, returnedTask.Title)
	assert.Equal(t, updatedTask.Status, returnedTask.Status)
	assert.Equal(t, createdTask.ID, returnedTask.ID)

	// --- Очистка тестовых данных ---
	defer func() {
		_, err := db.Exec("DELETE FROM tasks WHERE id=$1", createdTask.ID)
		if err != nil {
			t.Fatal("Error cleaning up test data:", err)
		}
	}()
}
