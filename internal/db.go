package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(host, user, pass, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, pass, dbname)
	return sql.Open("postgres", connStr)
}
