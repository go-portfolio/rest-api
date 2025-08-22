package main

import (
	"database/sql"
	"flag" // для чтения флагов командной строки
	"fmt"
	"log"

	"github.com/go-portfolio/rest-api/internal/config" // загрузка конфигурации приложения
	"github.com/go-portfolio/rest-api/internal/seed"
	"github.com/go-portfolio/rest-api/internal/server"   // HTTP-сервер и handler’ы
	"github.com/go-portfolio/rest-api/internal/services" // сервисы для работы с БД

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // драйвер базы для миграций
	_ "github.com/golang-migrate/migrate/v4/source/file"       // источник миграций — файлы
	_ "github.com/lib/pq"
)

// @title REST API Example
// @version 1.0
// @description REST API для пользователей и задач
// @host localhost:8080
// @BasePath /
func main() {
	// Флаг командной строки: если true, применяем миграции перед запуском сервера
	withMigrations := flag.Bool("with-migrations", false, "Применить все миграции до старта приложения")
	flag.Parse() // читаем флаги

	// Загружаем конфигурацию (например, из YAML + .env)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err) // если не удалось загрузить конфиг — завершаем приложение
	}

	// Получаем строку подключения к базе данных (DSN)
	dbURL := cfg.DSN()

	// Если указан флаг --with-migrations, применяем все миграции
	if *withMigrations {
		fmt.Println("Applying migrations...")
		if err := applyMigrations(dbURL, cfg.Migrations.Path); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations applied successfully!")
	}

	// Подключаемся к PostgreSQL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err) // завершаем приложение, если не удалось подключиться
	}

	seed.SeedUsers(db)
	// Создаём сервис для работы с задачами, используя реальную базу
	// Этот сервис реализует интерфейс TaskService
	taskSvc := services.NewPostgresTaskService(db)
	userSvc := services.NewPostgresUserService(db)

	fmt.Println("Starting application...")
	// Передаём сервис в сервер и запускаем HTTP-сервер
	server.StartServer(taskSvc, userSvc, cfg)
}

// applyMigrations применяет все миграции из указанной папки к базе данных
func applyMigrations(dbURL, migrationsPath string) error {
	// Создаём объект миграций
	m, err := migrate.New(
		"file://"+migrationsPath, // путь к папке с миграциями
		dbURL,                    // строка подключения к БД
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Применяем миграции "вверх" (Up)
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange { // ErrNoChange — если миграций нет
		return fmt.Errorf("migration up failed: %w", err)
	}

	return nil
}
