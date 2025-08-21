package services

import (
	"database/sql" // стандартная библиотека для работы с SQL-базами

	"github.com/go-portfolio/rest-api/internal/models" // структура Task
)

// -----------------------------
// Интерфейс TaskService
// -----------------------------
// Определяет методы для работы с задачами.
// Благодаря интерфейсу можно подставлять разные реализации:
// - реальную (PostgresTaskService)
// - мок для тестов (MockTaskService)
type TaskService interface {
	// Получить все задачи
	GetTasks() ([]models.Task, error)
	// Создать новую задачу и вернуть её ID
	CreateTask(title, status string) (int, error)
}

// -----------------------------
// Реализация TaskService для PostgreSQL
// -----------------------------
type PostgresTaskService struct {
	DB *sql.DB // подключение к базе данных
}

// Конструктор PostgresTaskService
func NewPostgresTaskService(db *sql.DB) *PostgresTaskService {
	// Возвращает указатель на новую структуру с подключением к БД
	return &PostgresTaskService{DB: db}
}

// -----------------------------
// Метод GetTasks
// -----------------------------
func (p *PostgresTaskService) GetTasks() ([]models.Task, error) {
	// Выполняем SQL-запрос для получения всех задач
	rows, err := p.DB.Query("SELECT id, title, status FROM tasks")
	if err != nil {
		// Если ошибка при запросе — возвращаем её
		return nil, err
	}
	// Обязательно закрываем rows после использования
	defer rows.Close()

	var tasks []models.Task

	// Проходим по всем строкам результата
	for rows.Next() {
		var t models.Task
		// Сканируем значения в структуру Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
			return nil, err
		}
		// Добавляем задачу в срез
		tasks = append(tasks, t)
	}

	// Возвращаем срез задач
	return tasks, nil
}

// -----------------------------
// Метод CreateTask
// -----------------------------
func (p *PostgresTaskService) CreateTask(title, status string) (int, error) {
	var id int
	// Выполняем INSERT и сразу возвращаем сгенерированный ID
	err := p.DB.QueryRow(
		"INSERT INTO tasks(title, status) VALUES($1, $2) RETURNING id",
		title, status,
	).Scan(&id) // сканируем результат (ID) в переменную

	if err != nil {
		// Если ошибка при вставке — возвращаем её
		return 0, err
	}

	// Возвращаем ID созданной задачи
	return id, nil
}
