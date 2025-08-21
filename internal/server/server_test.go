package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services"
)

// TestTasksHandler тестирует обработчик /tasks с использованием мок-сервиса
func TestTasksHandler(t *testing.T) {
	// Создаём мок-сервис, который реализует интерфейс TaskService
	// Он возвращает заранее определённые данные и не обращается к реальной БД
	mockSvc := &services.MockTaskService{}

	// -----------------------------
	// Тестируем GET /tasks
	// -----------------------------
	t.Run("GET /tasks", func(t *testing.T) {
		// Создаём HTTP-запрос GET на маршрут /tasks
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		// httptest.NewRecorder() — объект, который "ловит" ответ сервера для проверки
		w := httptest.NewRecorder()

		// Получаем handler для нашего мок-сервиса
		handler := TasksHandler(mockSvc)
		// Вызываем handler как реальный HTTP-запрос
		handler(w, req)

		// Получаем результат
		resp := w.Result()
		defer resp.Body.Close() // обязательно закрываем тело ответа

		// Проверяем статус код
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Декодируем JSON-ответ в срез задач
		var tasks []models.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Fatal(err) // если не удалось распарсить JSON — тест падает
		}

		// Проверяем содержимое: должно быть ровно 1 задача с ID=1
		if len(tasks) != 1 || tasks[0].ID != 1 {
			t.Errorf("Unexpected tasks: %+v", tasks)
		}
	})

	// -----------------------------
	// Тестируем POST /tasks
	// -----------------------------
	t.Run("POST /tasks", func(t *testing.T) {
		// Создаём новую задачу для отправки в теле запроса
		newTask := models.Task{Title: "New Task", Status: "Open"}
		// Кодируем её в JSON
		body, _ := json.Marshal(newTask)

		// Создаём HTTP-запрос POST на /tasks с телом body
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
		// httptest.NewRecorder() — объект для записи ответа
		w := httptest.NewRecorder()

		// Получаем handler для мок-сервиса
		handler := TasksHandler(mockSvc)
		// Вызываем handler
		handler(w, req)

		// Получаем результат
		resp := w.Result()
		defer resp.Body.Close()

		// Проверяем статус код
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Декодируем ответ в структуру Task
		var created models.Task
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatal(err)
		}

		// Проверяем, что созданная задача имеет ожидаемый ID и Title
		if created.ID != 42 || created.Title != newTask.Title {
			t.Errorf("Unexpected created task: %+v", created)
		}
	})
}
