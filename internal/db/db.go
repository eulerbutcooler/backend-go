package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	connStr := fmt.Sprintf("host=db user=postgres password=secret dbname=bckndtest sslmode=disable")
	var err error
	maxRetries := 10
	for range maxRetries {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		if err := db.Ping(); err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
		db.Close()
	}
	return fmt.Errorf("failed to ping database after %d retries: %w", maxRetries, err)
}

func GetDB() *sql.DB {
	return db
}
