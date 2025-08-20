# Путь к файлу конфигурации проекта
CONFIG_PATH=configs/config.yaml

# Получение строки подключения к базе данных (DSN) через утилиту на Go
DB_URL=$(shell go run ./cmd/utils/config_printer.go dsn)

# Получение пути к миграциям через утилиту на Go
MIGRATIONS_PATH=$(shell go run ./cmd/utils/config_printer.go migrations_path)

# TODO: добавить реализацию запуска всех миграций
# Запуск приложения с выполнением миграций
#run:
#	go run ./cmd/app/main.go --with-migrations

# Применение всех миграций к базе данных
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Откат миграций на один шаг вниз
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

# Принудительная установка версии миграций на 1
# Полезно, если база данных в "грязном" состоянии и нужно синхронизировать версию
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)
