package config

import (
	"fmt"
	"os"
	"path/filepath"
	"log"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config описывает структуру конфигурации приложения
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SslMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Migrations struct {
		Path string `yaml:"path"`
	} `yaml:"migrations"`
}

// LoadConfig загружает конфигурацию из YAML и .env
func LoadConfig() (*Config, error) {
	configPath := сonfigPath()
	configDir := filepath.Dir(configPath)
	envPath := filepath.Join(configDir, "..", ".env")
	_ = godotenv.Load(envPath)

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	// Подстановка переменных окружения
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		cfg.Database.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		cfg.Database.SslMode = v
	}

	return cfg, nil
}

// DSN возвращает готовую строку подключения к базе
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SslMode,
	)
}



// ConfigPath автоматически ищет configs/config.yaml
func сonfigPath() string {
	// 1. Сначала смотрим переменную окружения CONFIG_PATH
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// 2. Начинаем поиск от текущей рабочей директории
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get working directory: %v", err)
	}

	// 3. Поднимаемся вверх по дереву директорий, пока не найдём папку configs/config.yaml
	dir := wd
	for {
		tryPath := filepath.Join(dir, "configs", "config.yaml")
		if _, err := os.Stat(tryPath); err == nil {
			return tryPath
		}

		// Поднимаемся на уровень выше
		parent := filepath.Dir(dir)
		if parent == dir { // достигли корня FS
			break
		}
		dir = parent
	}

	log.Fatal("Cannot find configs/config.yaml in any parent directory")
	return ""
}