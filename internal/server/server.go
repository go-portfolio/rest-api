package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services"
)

// TasksHandler создаёт HTTP-обработчик для маршрута /tasks
// svc — интерфейс TaskService, через который handler взаимодействует с задачами
// Возвращает http.HandlerFunc, что позволяет передавать его в mux.HandleFunc
func TasksHandler(svc services.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Определяем поведение в зависимости от HTTP-метода
		switch r.Method {

		// -----------------------------
		// Обработка GET /tasks
		// -----------------------------
		case http.MethodGet:
			// Получаем список задач через сервис
			tasks, err := svc.GetTasks()
			if err != nil {
				// Если произошла ошибка — возвращаем 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Кодируем список задач в JSON и отправляем в ответ
			json.NewEncoder(w).Encode(tasks)

		// -----------------------------
		// Обработка POST /tasks
		// -----------------------------
		case http.MethodPost:
			// Декодируем JSON-тело запроса в структуру Task
			var t models.Task
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				// Если JSON некорректен — возвращаем 400 Bad Request
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Создаём новую задачу через сервис, получаем её ID
			id, err := svc.CreateTask(t.Title, t.Status)
			if err != nil {
				// Если ошибка при создании в сервисе — возвращаем 500
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Устанавливаем ID созданной задачи
			t.ID = id
			// Отправляем созданную задачу обратно клиенту в формате JSON
			json.NewEncoder(w).Encode(t)

		// -----------------------------
		// Если метод не GET и не POST
		// -----------------------------
		default:
			// Возвращаем 405 Method Not Allowed для остальных методов
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// StartServer запускает HTTP-сервер на порту 8080
// svc — интерфейс TaskService, чтобы обработчики могли работать с задачами
func StartServer(svc services.TaskService) {
	// Создаём новый HTTP-мультиплексор (router)
	mux := http.NewServeMux()
	// Регистрируем маршрут /tasks и привязываем к нему handler
	mux.HandleFunc("/tasks", TasksHandler(svc))
	// Запускаем HTTP-сервер на порту 8080
	// В реальном приложении можно добавить логирование и graceful shutdown
	http.ListenAndServe(":8080", mux)
}
