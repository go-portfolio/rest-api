package server

import (
	"encoding/json"

	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-portfolio/rest-api/docs" // docs генерируется swag
	"github.com/go-portfolio/rest-api/internal/auth"
	"github.com/go-portfolio/rest-api/internal/config"
	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
)

var requestCount = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    },
    []string{"path", "method"},
)

func init() {
    prometheus.MustRegister(requestCount)
}

var taskValidate = validator.New()

// TasksHandler godoc
// @Summary      Управление задачами
// @Description  Получение, создание, обновление и удаление задач
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id      path      int          false  "ID задачи"  example(1)
// @Param        task    body      models.Task  false  "Данные задачи"
// @Success      200     {array}   models.Task        "Список задач или обновленная задача"
// @Success      201     {object}  models.Task        "Созданная задача"
// @Success      204     {string}  string             "Задача удалена"
// @Failure      400     {string}  string             "Некорректный запрос"
// @Failure      401     {string}  string             "Неавторизован"
// @Failure      404     {string}  string             "Задача не найдена"
// @Failure      500     {string}  string             "Внутренняя ошибка сервера"
// @Router       /tasks [get]
// @Router       /tasks [post]
// @Router       /tasks/{id} [put]
// @Router       /tasks/{id} [delete]
func TasksHandler(svc services.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Считаем количество запросов к /tasks:
		requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()

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
			// простая"валидация" query-параметров
			limitStr := r.URL.Query().Get("limit")
			if limitStr != "" {
				limit, err := strconv.Atoi(limitStr)
				if err != nil || limit < 1 {
					http.Error(w, "invalid limit", http.StatusBadRequest)
					return
				}
				// используем limit
			}

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
			var t models.Task

			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if err := taskValidate.Struct(t); err != nil {
				errors := make(map[string]string)
				for _, e := range err.(validator.ValidationErrors) {
					errors[e.Field()] = e.Tag()
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errors)
				return
			}

			// Создаём новую задачу через сервис, получаем её ID
			id, err := svc.CreateTask(t.UserID, t.Title, t.Status)
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
			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				http.Error(w, "invalid task ID", http.StatusBadRequest)
				return
			}

			// 2. Декодируем тело запроса
			var t models.Task
			if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			// 3. Валидируем JSON
			if err := taskValidate.Struct(t); err != nil {
				errors := make(map[string]string)
				for _, e := range err.(validator.ValidationErrors) {
					errors[e.Field()] = e.Tag()
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errors)
				return
			}
			updated, err := svc.UpdateTask(taskID, t.UserID, t.Title, t.Status)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(updated)

		// -----------------------------
		// DELETE /tasks/{id}
		// -----------------------------
		case http.MethodDelete:
			// Вытаскиваем ID из пути
			idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				http.Error(w, "invalid task ID", http.StatusBadRequest)
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
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	// Новый роут для метрик
    mux.Handle("/metrics", promhttp.Handler())

	// Запускаем HTTP-сервер на порту 8080
	// В реальном приложении можно добавить логирование и graceful shutdown
	http.ListenAndServe(":8080", mux)
}
