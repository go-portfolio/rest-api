# Получаем пути и DSN через Go утилиту
MIGRATIONS_PATH=$(shell go run ./cmd/utils/config_printer.go migrations_path)
DB_URL=$(shell go run ./cmd/utils/config_printer.go dsn)

# Запуск приложения
run:
	go run ./cmd/app/main.go --with-migrations

# Применить все миграции
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Откат всех миграций
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

# Принудительная установка версии миграций
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(or $(VERSION),1)

# Сброс и повторное применение всех миграций
migrate-reset:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" drop -f
	make --no-print-directory migrate-up

