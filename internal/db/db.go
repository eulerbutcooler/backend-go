package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenvOrDefault("DB_HOST", "db"),
		getenvOrDefault("DB_PORT", "5432"),
		getenvOrDefault("DB_USER", "postgres"),
		getenvOrDefault("DB_PASSWORD", "secret"),
		getenvOrDefault("DB_NAME", "backend-db"),
		getenvOrDefault("DB_SSLMODE", "disable"),
	)
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

func getenvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDB() *sql.DB {
	return db
}
