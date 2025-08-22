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

# Схема работы TasksHandler с сервисом и моками
+--------+          +----------------+          +----------------------+
| Client |  GET/POST | HTTP Handler  | TaskService |  Реальная БД / Mock  |
+--------+  ------>  +----------------+ ------>  +----------------------+
              /tasks          |                        |
                             v                        v
                      TasksHandler(svc)        svc.GetTasks() / svc.CreateTask()
                             |                        |
                             |                        |
                    JSON Encode/Decode           Возврат данных или фиктивных ID
                             |
                             v
                       Response (JSON)
                             |
                             v
                          Client

🔹 Объяснение шагов

## Client

Отправляет HTTP-запрос (GET /tasks или POST /tasks).

HTTP Handler (TasksHandler)

Получает интерфейс TaskService (svc).

В зависимости от метода (GET или POST) вызывает соответствующий метод сервиса.

Кодирует результат в JSON и отправляет обратно клиенту.

## TaskService

Может быть:

Реальная реализация (PostgresTaskService) — работает с базой данных.

Mock (MockTaskService) — возвращает фиктивные данные для юнит-тестов.

## Реальная БД / Mock

В случае реальной БД выполняются SQL-запросы (SELECT или INSERT).

В случае мока возвращаются заранее определённые данные или фиксированные ID.

## Response

Handler кодирует результат в JSON и возвращает клиенту.

Клиент получает ответ без знания, используется ли мок или реальная база.

## Работа с API

### 1. Получение JWT токена через `/login`

Отправляем POST-запрос с логином и паролем пользователя:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alex","password":"password123"}'
```  
Пример ответа:
```bash
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU4NTAzNDMsInVzZXJfaWQiOjF9.Oi05RSu-G1_y9CbfIQ7NYFkmpJYxoX7cCA-kYrQGEGw"
}
```
token — это ваш JWT, который нужно использовать для авторизации при последующих запросах.

2. Получение списка задач через /tasks
Теперь можно получить задачи, передав токен в заголовке Authorization:
```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer <ваш_JWT_токен>"
```  
Пример с подставленным токеном:
```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU4NTAzNDMsInVzZXJfaWQiOjF9.Oi05RSu-G1_y9CbfIQ7NYFkmpJYxoX7cCA-kYrQGEGw"
```  
Пример ответа:
```json
[
  {"id":1,"title":"Test task","status":"todo"},
  {"id":2,"title":"Test task","status":"todo"}
]
```
