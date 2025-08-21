package models

type Task struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Status string `json:"status"`
}