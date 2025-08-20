package main

import (
	"fmt"
	"log"
	"os"
	"gopkg.in/yaml.v3"
)

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

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	// Подставляем переменные окружения, если они есть
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

	return cfg, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Specify argument: 'dsn' or 'migrations_path'")
	}

	cfg, err := LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "dsn":
		fmt.Printf("postgres://%s:%s@%s:%s/%s?sslmode=%s\n",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SslMode,
		)
	case "migrations_path":
		fmt.Println(cfg.Migrations.Path)
	default:
		log.Fatal("Unknown argument: use 'dsn' or 'migrations_path'")
	}
}
