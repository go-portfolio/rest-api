package services

import (
	"database/sql" // стандартная библиотека для работы с SQL-базами
	"time"

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
	CreateTask(userID int, title, status string) (int, error)
	UpdateTask(id int, userID int, title, status string) (*models.Task, error)
	DeleteTask(id int) error
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
	rows, err := p.DB.Query("SELECT id, title, status FROM tasks  WHERE deleted_at IS NULL")
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
func (p *PostgresTaskService) CreateTask(userID int, title, status string) (int, error) {
	var id int
	// Выполняем INSERT и сразу возвращаем сгенерированный ID
	err := p.DB.QueryRow(
		"INSERT INTO tasks(title, status, created_at, user_id) VALUES($1, $2, NOW(), $3) RETURNING id",
		title, status, userID,
	).Scan(&id) // сканируем результат (ID) в переменную

	if err != nil {
		// Если ошибка при вставке — возвращаем её
		return 0, err
	}

	// Возвращаем ID созданной задачи
	return id, nil
}

func (s *PostgresTaskService) UpdateTask(id int, userID int, title, status string) (*models.Task, error) {
	_, err := s.DB.Exec(`UPDATE tasks SET title=$1, status=$2, updated_at=NOW(), user_id=$4 WHERE id=$3 AND deleted_at IS NULL`, title, status, id, userID)
	if err != nil {
		return nil, err
	}

	var t models.Task
	err = s.DB.QueryRow(`SELECT id, title, status, user_id FROM tasks WHERE id=$1`, id).Scan(&t.ID, &t.Title, &t.Status, &t.UserID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *PostgresTaskService) DeleteTask(id int) error {
	now := time.Now()
	_, err := s.DB.Exec("UPDATE tasks SET deleted_at=$1 WHERE id=$2", now, id)
	return err
}
