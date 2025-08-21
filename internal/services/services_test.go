package services

import (
    "database/sql"
    "fmt"
    "testing"

    "github.com/go-portfolio/rest-api/internal/config" 
    _ "github.com/lib/pq"
)

func TestGetTasks(t *testing.T) {
    // Загружаем конфиг из корня проекта
    cfg, err := config.LoadConfig("../../configs/config.yaml")
    if err != nil {
        t.Fatal(err)
    }

    // Формируем DSN
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Name,
        cfg.Database.SslMode,
    )

    // Подключаемся к БД
    db, err := sql.Open("postgres", dsn)
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
