package services

import (
    "database/sql"
    "testing"

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

    // Проверяем сервис
    tasks, err := GetTasks(db)
    if err != nil {
        t.Fatal(err)
    }

    if len(tasks) == 0 {
        t.Error("expected at least one task")
    }
}
