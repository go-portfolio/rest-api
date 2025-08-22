package unit

import (
    "database/sql"
    "testing"
    "log"

    "github.com/go-portfolio/rest-api/internal/config" 
    _ "github.com/lib/pq" 
	"github.com/go-portfolio/rest-api/internal/services" 
)

// ------------------------
// Интеграционный тест: получение задач из БД
// ------------------------
func TestGetTasks(t *testing.T) {
    // Загружаем конфиг из корня проекта (DSN, настройки БД и т.п.)
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err) // Завершаем тест при ошибке загрузки конфига
    }

    // Подключаемся к базе данных PostgreSQL
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close() // Закрываем соединение после теста

    // Создаем сервис для работы с задачами
    svc := services.NewPostgresTaskService(db)
    
    // Получаем список задач
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatal(err)
    }

    // Проверяем, что список задач не пустой
    if len(tasks) == 0 {
        t.Error("expected at least one task")
    }
}

// ------------------------
// Интеграционный тест: создание и проверка задачи
// ------------------------
func TestCreateAndGetTasks(t *testing.T) {
    // Загружаем конфиг проекта
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err)
    }

    // Подключаемся к БД
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Создаем сервис для работы с задачами
    svc := services.NewPostgresTaskService(db)

    // Создаем тестовую задачу
    id, err := svc.CreateTask(1, "Test task", "todo")
    if err != nil {
        t.Fatalf("Failed to create task: %v", err)
    }

    // Получаем список задач после создания
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatalf("Failed to get tasks: %v", err)
    }

    // Проверяем, что созданная задача присутствует в списке
    found := false
    for _, task := range tasks {
        if task.ID == id {
            found = true
            break
        }
    }

    if !found {
        t.Errorf("Таск не найден после создания")
    }
}

// ------------------------
// Интеграционный тест: обновление задачи
// ------------------------
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

    svc := services.NewPostgresTaskService(db)

    // 1. Создаем задачу
    id, err := svc.CreateTask(1, "Initial task", "todo")
    if err != nil {
        t.Fatalf("Failed to create task: %v", err)
    }

    // 2. Обновляем задачу
    _, err = svc.UpdateTask(id, 1, "Updated task", "in-progress")
    if err != nil {
        t.Fatalf("Failed to update task: %v", err)
    }

    // 3. Получаем список задач и проверяем обновление
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatalf("Failed to get tasks: %v", err)
    }

    found := false
    for _, task := range tasks {
        if task.ID == id {
            found = true
            if task.Title != "Updated task" || task.Status != "in-progress" {
                t.Errorf("Task not updated correctly: got %+v", task)
            }
        }
    }

    if !found {
        t.Errorf("Task with ID %d not found after update", id)
    }
}

// ------------------------
// Интеграционный тест: удаление задачи
// ------------------------
func TestDeleteTask(t *testing.T) {
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err)
    }

    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    svc := services.NewPostgresTaskService(db)

    // 1. Создаем задачу
    id, err := svc.CreateTask(1, "Task to delete", "todo")
    if err != nil {
        t.Fatalf("Failed to create task: %v", err)
    }

    // 2. Удаляем задачу
    err = svc.DeleteTask(id)
    if err != nil {
        t.Fatalf("Failed to delete task: %v", err)
    }

    // 3. Проверяем, что задачи нет в списке
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatalf("Failed to get tasks: %v", err)
    }

    for _, task := range tasks {
        if task.ID == id {
            t.Errorf("Task with ID %d still exists after deletion", id)
        }
    }
}
