package services

import (
	"errors"

	"github.com/go-portfolio/rest-api/internal/models"
)

// -----------------------------
// MockTaskService
// -----------------------------
// Мок-реализация интерфейса TaskService для юнит-тестов.
// Позволяет тестировать обработчики и другие компоненты
// без реального подключения к базе данных.
type MockTaskService struct{
	Tasks []models.Task
}

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

// -----------------------------
// UpdateTask
// -----------------------------
func (m *MockTaskService) UpdateTask(id int, title, status string) (*models.Task, error) {
	for i, t := range m.Tasks {
		if t.ID == id {
			m.Tasks[i].Title = title
			m.Tasks[i].Status = status
			return &m.Tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}

// -----------------------------
// DeleteTask
// -----------------------------
func (m *MockTaskService) DeleteTask(id int) error {
	for i, t := range m.Tasks {
		if t.ID == id {
			m.Tasks = append(m.Tasks[:i], m.Tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
