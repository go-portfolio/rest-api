package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/go-portfolio/rest-api/internal/config"
	"github.com/go-portfolio/rest-api/internal/server"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // драйвер базы для миграций
	_ "github.com/golang-migrate/migrate/v4/source/file"       // источник миграций — файлы
	_ "github.com/lib/pq"                                       // драйвер PostgreSQL для sql.Open
)

func main() {
	// Флаг для применения миграций перед запуском приложения
	withMigrations := flag.Bool("with-migrations", false, "Apply database migrations before starting the app")
	flag.Parse()

	// Загружаем конфигурацию приложения (YAML + .env)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err) // если конфиг не загрузился — выходим
	}

	// Получаем готовую строку подключения к базе данных (DSN)
	dbURL := cfg.DSN()

	// Если указан флаг с миграциями, применяем их
	if *withMigrations {
		fmt.Println("Applying migrations...")
		if err := applyMigrations(dbURL, cfg.Migrations.Path); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations applied successfully!")
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting application...")
	// Запускаем HTTP-сервер с маршрутами, используя подключение к базе
	server.StartServer(db)
}

// applyMigrations применяет все миграции из указанной папки к базе данных
func applyMigrations(dbURL, migrationsPath string) error {
	// Создаем экземпляр миграций
	m, err := migrate.New(
		"file://"+migrationsPath, // путь к миграциям
		dbURL,                     // строка подключения к БД
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Применяем миграции
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange { // ErrNoChange — когда миграций нет
		return fmt.Errorf("migration up failed: %w", err)
	}

	return nil
}
