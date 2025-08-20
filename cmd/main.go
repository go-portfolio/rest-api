package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

// Хранит настройки для подключения к БД
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func main() {
	fmt.Println("REST API запускается...")
	// Открываем YAML файл
	f, err := os.Open("../configs/config.yaml")
	if err != nil {
		log.Fatal("Не удалось открыть config.yaml:", err)
	}
	defer f.Close()

	// Читаем и парсим YAML
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatal("Ошибка парсинга config.yaml:", err)
	}

	// Формируем строку подключения к PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	// Подключаемся к базе
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка при открытии подключения:", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к PostgreSQL:", err)
	}

	fmt.Println("Подключение к PostgreSQL успешно!")
}
