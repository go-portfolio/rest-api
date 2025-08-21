package services

import (
    "database/sql"
    "github.com/go-portfolio/rest-api/internal/models"
)

func GetTasks(db *sql.DB) ([]models.Task, error) {
    rows, err := db.Query("SELECT id, title, status FROM tasks")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []models.Task

    for rows.Next() {
        var t models.Task
        if err := rows.Scan(&t.ID, &t.Title, &t.Status); err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }

    return tasks, nil
}
