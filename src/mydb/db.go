package mydb

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"main/src/logger" // Используйте ваш логгер
)

// setenv загружает переменные окружения и возвращает строку подключения.
func setenv() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
	)
	logger.Logger.Debug("Connecting to database with connection string:", dsn)
	return dsn, nil
}

// Database подключается к базе данных и возвращает *sql.DB.
func Database() (*sql.DB, error) {
	dsn, err := setenv()
	if err != nil {
		logger.Logger.Error("Failed to load .env file", "error", err)
		return nil, err
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Logger.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Logger.Error("Failed to ping database", "error", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Logger.Info("Successfully connected to database")
	return db, nil
}