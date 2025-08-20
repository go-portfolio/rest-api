package config

import (
	"log"  // для логирования ошибок
	"os"   // для чтения файлов
	"fmt"
	"gopkg.in/yaml.v3" // библиотека для работы с YAML
)

// Config — структура для хранения всей конфигурации приложения
type Config struct {
	// Database — настройки подключения к PostgreSQL
	Database struct {
		User     string `yaml:"user"`     // имя пользователя базы
		Password string `yaml:"password"` // пароль пользователя
		Host     string `yaml:"host"`     // адрес сервера базы
		Port     int    `yaml:"port"`     // порт базы данных
		Name     string `yaml:"name"`     // имя базы данных
		SslMode  string `yaml:"sslmode"`  // sslmode (disable, require и т.д.)
	} `yaml:"database"`

	// Migrations — путь к файлам миграций
	Migrations struct {
		Path string `yaml:"path"` // директория с SQL миграциями
	} `yaml:"migrations"`
}

// LoadConfig — функция для загрузки конфигурации из YAML файла
func LoadConfig(path string) *Config {
	// Читаем содержимое файла
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	var cfg Config
	// Парсим YAML в структуру Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}

	return &cfg // возвращаем указатель на структуру
}

// DBUrl — формирует строку подключения к PostgreSQL в формате DSN
func (c *Config) DBUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,   // вот здесь число останется числом
		c.Database.Name,
		c.Database.SslMode,
	)
}


