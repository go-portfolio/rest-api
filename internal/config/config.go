package config

import (
    "os"
    "github.com/joho/godotenv"
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
    // Подгружаем .env локально
    godotenv.Load()

    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    cfg := &Config{}
    decoder := yaml.NewDecoder(f)
    if err := decoder.Decode(cfg); err != nil {
        return nil, err
    }

    // Заменяем переменные окружения
    cfg.Database.Host = os.Getenv("DB_HOST")
    cfg.Database.Port = os.Getenv("DB_PORT")
    cfg.Database.User = os.Getenv("DB_USER")
    cfg.Database.Password = os.Getenv("DB_PASSWORD")
    cfg.Database.Name = os.Getenv("DB_NAME")

    return cfg, nil
}
