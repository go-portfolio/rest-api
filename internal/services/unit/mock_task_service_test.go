package unit

import (
	"testing"

	"github.com/go-portfolio/rest-api/internal/models"
	"github.com/go-portfolio/rest-api/internal/services"
)

// -----------------------------
// Тестирование GetTasks()
// -----------------------------
func TestMockTaskService_GetTasks(t *testing.T) {
	// Создаём мок-сервис
	mock := &services.MockTaskService{}

	// Вызываем GetTasks и проверяем результат
	tasks, err := mock.GetTasks()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем, что вернулся ровно один таск
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	// Проверяем, что заголовок задачи совпадает с ожидаемым
	if tasks[0].Title != "Test Task" {
		t.Errorf("expected task title 'Test Task', got %s", tasks[0].Title)
	}
}

// -----------------------------
// Тестирование CreateTask()
// -----------------------------
func TestMockTaskService_CreateTask(t *testing.T) {
	mock := &services.MockTaskService{}

	// Создаём задачу и проверяем возвращаемый ID
	id, err := mock.CreateTask(1, "New Task", "New")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id != 42 { // мок всегда возвращает 42
		t.Errorf("expected ID 42, got %d", id)
	}
}

// -----------------------------
// Тестирование UpdateTask()
// -----------------------------
func TestMockTaskService_UpdateTask(t *testing.T) {
	mock := &services.MockTaskService{
		Tasks: []models.Task{
			{ID: 1, Title: "Old Title", Status: "New"},
		},
	}

	// Обновляем существующую задачу
	task, err := mock.UpdateTask(1, 1, "Updated Title", "Done")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем, что данные обновились
	if task.Title != "Updated Title" || task.Status != "Done" {
		t.Errorf("task not updated correctly: %+v", task)
	}

	// Пробуем обновить несуществующую задачу
	_, err = mock.UpdateTask(99, 1, "X", "Y")
	if err == nil {
		t.Errorf("expected error for non-existent task, got nil")
	}
}

// -----------------------------
// Тестирование DeleteTask()
// -----------------------------
func TestMockTaskService_DeleteTask(t *testing.T) {
	mock := &services.MockTaskService{
		Tasks: []models.Task{
			{ID: 1, Title: "Task 1", Status: "New"},
			{ID: 2, Title: "Task 2", Status: "Done"},
		},
	}

	// Удаляем существующую задачу
	err := mock.DeleteTask(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем, что осталась только вторая задача
	if len(mock.Tasks) != 1 || mock.Tasks[0].ID != 2 {
		t.Errorf("task not deleted correctly, remaining tasks: %+v", mock.Tasks)
	}

	// Пробуем удалить несуществующую задачу
	err = mock.DeleteTask(99)
	if err == nil {
		t.Errorf("expected error for non-existent task, got nil")
	}
}
