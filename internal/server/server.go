package server

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "github.com/go-portfolio/rest-api/internal/services"
)

func StartServer(db *sql.DB) {
    mux := http.NewServeMux()
    mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        tasks, err := services.GetTasks(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(tasks)
    })

    http.ListenAndServe(":8080", mux)
}
