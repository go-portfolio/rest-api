package services

import (
    "database/sql"
    "testing"
    "log"

    "github.com/go-portfolio/rest-api/internal/config" 
    _ "github.com/lib/pq"
)

func TestGetTasks(t *testing.T) {
    // Загружаем конфиг из корня проекта
    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatal(err)
    }

    // Подключаемся к БД
    db, err := sql.Open("postgres", cfg.DSN())
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

	// Создаем сервис
	svc := NewPostgresTaskService(db)
	
    // Проверяем сервис
    tasks, err := svc.GetTasks()
    if err != nil {
        t.Fatal(err)
    }

    if len(tasks) == 0 {
        t.Error("expected at least one task")
    }
}

func TestCreateAndGetTasks(t *testing.T) {
	// Загружаем конфиг из корня проекта
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

	// Создаем сервис
	svc := NewPostgresTaskService(db)

	// Создаем тестовую задачу
	id, err := svc.CreateTask("Test task", "todo")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	tasks, err := svc.GetTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

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
