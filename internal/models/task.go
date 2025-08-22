package models

import "time"

// Task представляет задачу
// swagger:model Task
type Task struct {
	// ID задачи
    // example: 1
	ID        int       `json:"id"`
	// Заголовок задачи
    // example: "Сделать домашку"
	Title     string    `json:"title"`
	// Статус задачи
    // example: "pending"
	Status    string    `json:"status"`
	// Дата создания задачи в формате RFC3339
    // example: "2025-08-22T17:00:00Z"
	CreatedAt time.Time `json:"created_at"`
	// Дата обновления задачи в формате RFC3339
    // example: "2025-08-22T17:00:00Z"
	UpdatedAt time.Time `json:"updated_at"`
	// Дата soft-удаления задачи в формате RFC3339
    // example: "2025-08-22T17:00:00Z"
	DeletedAt *time.Time `db:"deleted_at"` // soft-удаления
	// ID пользователя, которому принадлежит задача
    // example: 42
    UserID    int       `json:"user_id" db:"user_id"`
}
