# REST API Project

---

## Структура проекта

rest-api/
├─ cmd/
│ ├─ app/ # точка входа приложения
│ └─ utils/ # утилиты, например config_printer.go
├─ internal/
│ └─ config/ # чтение конфигураций
├─ migrations/ # SQL миграции
├─ configs/
│ └─ config.yaml # шаблон конфигурации
├─ Makefile
└─ .env.example # пример переменных окружения


---

## Установка и запуск

### 1. Создать `.env`  

Скопируйте пример и заполните локальные данные:

```bash
cp .env.example .env
Пример .env:

DB_HOST=localhost
DB_PORT=5432
DB_USER=alex
DB_PASSWORD=secret
DB_NAME=restapi
Важно: .env не коммитится и хранит секреты локально.

2. Установить зависимости Go
go mod tidy
3. Makefile — основные команды
Команда	Описание
make run	Запуск приложения с применением миграций
make migrate-up	Применить все новые миграции
make migrate-down	Откатить только последнюю миграцию
make migrate-force	Принудительно установить версию миграций. Можно указать версию: make migrate-force VERSION=2
make migrate-reset	Полный сброс базы и повторное применение всех миграций

4. Пример работы
# Применить все миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Принудительно синхронизировать версию миграций
make migrate-force VERSION=1

# Сбросить базу и заново применить миграции
make migrate-reset
5. Безопасность и портфолио
Конфигурация и секреты не хранятся в репозитории.

Makefile использует .env и config_printer.go, чтобы безопасно подставлять DSN и путь к миграциям.

Проект можно запускать локально на любой машине, не раскрывая реальные пароли.

6. Примечания
Миграции создаются в папке migrations/ с префиксом номера версии:
001_create_tasks_table.up.sql
001_create_tasks_table.down.sql
Go-код использует internal/config и config_printer.go для универсального чтения DSN и пути миграций.

