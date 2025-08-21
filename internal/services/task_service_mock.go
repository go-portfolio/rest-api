package services

import "github.com/go-portfolio/rest-api/internal/models"

// -----------------------------
// MockTaskService
// -----------------------------
// Мок-реализация интерфейса TaskService для юнит-тестов.
// Позволяет тестировать обработчики и другие компоненты
// без реального подключения к базе данных.
type MockTaskService struct{}

// -----------------------------
// GetTasks
// -----------------------------
// Возвращает заранее заданный срез задач.
// Не обращается к реальной базе, просто имитирует результат.
// Возвращаемые данные позволяют проверить, что обработчик правильно декодирует JSON и возвращает список.
func (m *MockTaskService) GetTasks() ([]models.Task, error) {
	return []models.Task{
		{ID: 1, Title: "Test Task", Status: "New"},
	}, nil
}

// -----------------------------
// CreateTask
// -----------------------------
// Имитирует создание задачи и возвращает фиктивный ID (42).
// Не записывает данные в базу, позволяет проверить работу POST /tasks в тестах.
func (m *MockTaskService) CreateTask(title, status string) (int, error) {
	return 42, nil
}
