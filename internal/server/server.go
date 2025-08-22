package server

import (
	"encoding/json"

	"net/http"
	"strings"
	"strconv"

	"github.com/go-portfolio/rest-api/internal/auth"
	"github.com/go-portfolio/rest-api/internal/config"
	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services"
)

// TasksHandler создаёт HTTP-обработчик для маршрута /tasks
// svc — интерфейс TaskService, через который handler взаимодействует с задачами
// Возвращает http.HandlerFunc, что позволяет передавать его в mux.HandleFunc
func TasksHandler(svc services.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если URL содержит ID задачи (например, /tasks/1), извлекаем его
		pathParts := strings.Split(r.URL.Path, "/")
		var taskID int
		var err error
		if len(pathParts) == 3 && pathParts[2] != "" {
			taskID, err = strconv.Atoi(pathParts[2])
			if err != nil {
				http.Error(w, "invalid task id", http.StatusBadRequest)
				return
			}
		}

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
		// PUT /tasks/{id}
		// -----------------------------
		case http.MethodPut:
			if taskID == 0 {
				http.Error(w, "task id required", http.StatusBadRequest)
				return
			}
			var t models.Task
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updated, err := svc.UpdateTask(taskID, t.Title, t.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(updated)

		// -----------------------------
		// DELETE /tasks/{id}
		// -----------------------------
		case http.MethodDelete:
			if taskID == 0 {
				http.Error(w, "task id required", http.StatusBadRequest)
				return
			}
			if err := svc.DeleteTask(taskID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)	

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
func StartServer(svc services.TaskService, userSvc services.UserService, cfg *config.Config) {
	// Создаём новый HTTP-мультиплексор (router)
	mux := http.NewServeMux()
	// Public endpoints
	mux.HandleFunc("/login", LoginHandler(userSvc, cfg.Jwt.JwtSecretKey))
	// Регистрируем маршрут /tasks и привязываем к нему handler
	mux.Handle("/tasks", auth.VerifyToken(cfg.Jwt.JwtSecretKey)(TasksHandler(svc)))
	mux.Handle("/tasks/", auth.VerifyToken(cfg.Jwt.JwtSecretKey)(TasksHandler(svc)))

	// Запускаем HTTP-сервер на порту 8080
	// В реальном приложении можно добавить логирование и graceful shutdown
	http.ListenAndServe(":8080", mux)
}
