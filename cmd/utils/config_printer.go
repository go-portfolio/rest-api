package main

import (
    "fmt"
    "os"
    "github.com/go-portfolio/rest-api/internal/config"
)


func main() {
	cfg := config.LoadConfig("configs/config.yaml")

	if len(os.Args) < 2 {
		fmt.Println("expected arg: dsn | migrations_path")
		return
	}

	switch os.Args[1] {
	case "dsn":
		fmt.Println(cfg.DBUrl())
	case "migrations_path":
		fmt.Println(cfg.Migrations.Path)
	}
}
