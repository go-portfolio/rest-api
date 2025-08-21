package main

import (
	"fmt"                   // для вывода данных в консоль
	"log"                   // для логирования ошибок
	"os"                    // для работы с аргументами командной строки

	"github.com/go-portfolio/rest-api/internal/config" // пакет для работы с конфигурацией
)

func main() {
	// Проверяем, что пользователь передал аргумент командной строки
	if len(os.Args) < 2 {
		log.Fatal("Задайте аргументы: 'dsn' или 'migrations_path'") // завершаем программу с сообщением об ошибке
	}

	// Загружаем конфигурацию из файла config.yaml
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal(err) // завершаем программу, если не удалось загрузить конфигурацию
	}

	// В зависимости от переданного аргумента выполняем разные действия
	switch os.Args[1] {
	case "dsn":
		// Формируем строку подключения к базе данных PostgreSQL
		fmt.Printf("postgres://%s:%s@%s:%s/%s?sslmode=%s\n",
			cfg.Database.User,     // имя пользователя БД
			cfg.Database.Password, // пароль
			cfg.Database.Host,     // хост БД
			cfg.Database.Port,     // порт
			cfg.Database.Name,     // имя базы данных
			cfg.Database.SslMode,  // SSL режим
		)
	case "migrations_path":
		// Выводим путь к папке с миграциями
		fmt.Println(cfg.Migrations.Path)
	default:
		// Если аргумент неизвестен, завершаем программу с сообщением
		log.Fatal("Неизвестный аргумент: используй 'dsn' или 'migrations_path'")
	}
}
