package config

import (
	"os"
    "path/filepath"

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

// LoadConfig загружает конфигурацию из YAML и подставляет переменные окружения
func LoadConfig(configPath string) (*Config, error) {
    // configPath — путь к config.yaml
    configDir := filepath.Dir(configPath)
	// поднимаемся на уровень выше и загружаем .env
    envPath := filepath.Join(configDir, "..", ".env")
    _ = godotenv.Load(envPath)

	// Открываем YAML файл конфигурации
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	// Подставляем переменные окружения, если они установлены
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
