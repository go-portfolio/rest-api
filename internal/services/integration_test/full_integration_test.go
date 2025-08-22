package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-portfolio/rest-api/internal/auth"
	"github.com/go-portfolio/rest-api/internal/config"
	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/server"
	"github.com/go-portfolio/rest-api/internal/services"
	_ "github.com/lib/pq" // PostgreSQL драйвер
	"github.com/stretchr/testify/assert"
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

	var userID int
	// --------------------------
	// Создаём тестового пользователя в базе
	// Можно хранить пароль в явном виде для упрощения теста
	// --------------------------
	err = db.QueryRow(`INSERT INTO users (username, password_hash) VALUES ($1, $2)
                      ON CONFLICT (username) DO NOTHING RETURNING id`,
		"testuser", "password123").Scan(&userID) // Для реального проекта лучше использовать bcrypt
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Пользователь уже существует")
			// Можно, например, получить существующий ID:
			db.QueryRow("SELECT id FROM users WHERE username=$1", "testuser").Scan(&userID)
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Новый пользователь ID:", userID)
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
	newTask := models.Task{UserID: userID, Title: "Integration Task", Status: "todo"}
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

	// --------------------------
	// 4. Тестируем обновление задачи (PUT /tasks)
	// --------------------------
	updatedTask := models.Task{UserID: userID, Title: "Updated Task. Step 4", Status: "done"}
	updateBody, _ := json.Marshal(updatedTask)

	url := ts.URL + "/tasks/" + strconv.Itoa(createdTask.ID)
	updateReq, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := http.DefaultClient.Do(updateReq)
	if err != nil {
		t.Fatal(err)
	}
	defer updateResp.Body.Close()

	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	// Проверяем, что задача обновилась в базе
	row = db.QueryRow(`SELECT id, title, status FROM tasks WHERE id=$1`, createdTask.ID)
	err = row.Scan(&taskFromDB.ID, &taskFromDB.Title, &taskFromDB.Status)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, updatedTask.Title, taskFromDB.Title)
	assert.Equal(t, updatedTask.Status, taskFromDB.Status)

	// --------------------------
	// 5. Тестируем удаление задачи (DELETE /tasks)
	// --------------------------
	deleteReq, _ := http.NewRequest(http.MethodDelete, ts.URL+"/tasks/"+strconv.Itoa(createdTask.ID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)

	deleteResp, err := http.DefaultClient.Do(deleteReq)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	// Проверяем, что задача удалена из базы
	row = db.QueryRow(`SELECT id FROM tasks WHERE id=$1`, createdTask.ID)
	err = row.Scan(&taskFromDB.ID)
	assert.Error(t, err) // Ожидаем ошибку, так как задача должна быть удалена
}
