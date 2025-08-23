# REST API Project

---

## Структура проекта
```
├── cmd
│   ├── main.go
│   └── utils
│       └── config_printer.go
├── configs
│   └── config.yaml
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── auth
│   │   └── jwt.go
│   ├── config
│   │   └── config.go
│   ├── models
│   │   ├── task.go
│   │   └── user.go
│   ├── repositories
│   │   └── user_repo.go
│   ├── seed
│   │   └── seed.go
│   ├── server
│   │   ├── auth_handlers.go
│   │   ├── server.go
│   │   └── server_test.go
│   └── services
│       ├── auth_service.go
│       ├── auth_service_mock.go
│       ├── integration_test
│       ├── task_service.go
│       ├── task_service_mock.go
│       └── unit
├── LICENSE
├── Makefile
├── migrations
│   ├── 002_create_users_table.down.sql
│   ├── 002_create_users_table.up.sql
│   ├── 003_create_tasks_table.down.sql
│   └── 003_create_tasks_table.up.sql
├── project_structure.txt
├── prometheus.yml
├── README.md
└── tests
```

### Описание модулей

- **cmd/server/** — точка входа, запускает HTTP-сервер.  
- **internal/server/** — настройка роутов, middleware, запуск API.  
- **internal/services/** — бизнес-логика (работа с БД, CRUD для задач).  
- **internal/models/** — структуры данных (например, `Task`, `User`).  
- **internal/config/** — конфигурация проекта (чтение `.env`, DSN).  
- **migrations/** — SQL-миграции для базы данных.  
- **static/** — статические файлы (CSS, JS, картинки).  
- **Makefile** — команды для запуска, тестов и миграций.  
- **go.mod** — список зависимостей Go.  

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
⚠️ Важно: .env не коммитится и хранит секреты локально.

### 2. Установить зависимости Go
```bash
go mod tidy
```
### 3. Makefile — основные команды
Команда	Описание
```bash
make run	Запуск приложения с применением миграций
make migrate-up	Применить все новые миграции
make migrate-down	Откатить только последнюю миграцию
make migrate-force	Принудительно установить версию миграций. Можно указать версию: make migrate-force VERSION=2
make migrate-reset	Полный сброс базы и повторное применение всех миграций
```

### 4. Пример работы
```bash
# Применить все миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Принудительно синхронизировать версию миграций
make migrate-force VERSION=1

# Сбросить базу и заново применить миграции
make migrate-reset
```
### 5. Безопасность и портфолио
Конфигурация и секреты не хранятся в репозитории.

`Makefile` использует `.env` и `config_printer.go`, чтобы безопасно подставлять `DSN` и путь к миграциям.

Проект можно запускать локально на любой машине, не раскрывая реальные пароли.

###6. Примечания
Миграции создаются в папке migrations/ с префиксом номера версии:

```SQL
001_create_tasks_table.up.sql
001_create_tasks_table.down.sql
Go-код использует internal/config и config_printer.go для универсального чтения DSN и пути миграций.
```

#### Объяснение шагов
##### Client
Отправляет `HTTP-запрос` (`GET` `/tasks` или `POST` `/tasks)`.

`HTTP Handler` (`TasksHandler`)

Получает интерфейс `TaskService (svc)`.

В зависимости от метода (`GET`/`POST`) вызывает соответствующий метод сервиса.

Кодирует результат в `JSON` и отправляет обратно клиенту.

###### TaskService

Реальная реализация (`PostgresTaskService`) — работает с базой данных.

`Mock` (`MockTaskService`) — возвращает фиктивные данные для юнит-тестов.

Реальная БД / `Mock`

В случае реальной БД выполняются `SQL`-запросы.

В случае мока возвращаются заранее определённые данные или фиксированные `ID`.

##### Response
`Handler` кодирует результат в `JSON` и возвращает клиенту.

## Работа с API
### 1. Получение `JWT`-токена через `/login`
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alex","password":"password123"}'
```  
Пример ответа:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR..."
}
```
token — это ваш `JWT`, который нужно использовать для авторизации при последующих запросах.

### 2. Получение списка задач через `/tasks`
```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer <ваш_JWT_токен>"
```
#### Пример ответа:

```json
[
  {"id":1,"title":"Test task","status":"todo"},
  {"id":2,"title":"Test task","status":"todo"}
]
```
#### Команды для запуска тестов
Unit-тесты
`Services`

```bash
# Запуск всех тестов
go test ./... -v -count=1

# Запуск всех unit-тестов
go test ./internal/services/unit -v -count=1

# Запуск конкретного unit-теста
go test ./internal/services/unit -v -run TestMockUserService_Authenticate
```
Server

```bash
# Запуск всех unit-тестов
go test ./internal/server -v -count=1

# Запуск конкретного unit-теста
go test ./internal/server -v -run TestTasksHandler
```
Интеграционные тесты
```bash
# Запуск всех интеграционных тестов
go test ./internal/services/integration_test -v -count=1

# Запуск конкретного интеграционного теста
go test ./internal/services/integration_test -v -run TestFullIntegration

# Просмотр доступных тестов
go test ./internal/services/integration_test -list .
```
## Grafana

Для запуска Grafana выполните команду:

```bash
docker compose up -d
```
После этого:

Откройте Grafana по ссылке: http://localhost:3000/

Логин: admin

Пароль: admin

Для просмотра Prometheus откройте: http://localhost:9090/
