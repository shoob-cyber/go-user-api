package database

import (
	"database/sql"

	"go-user-api/internal/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// Connect establishes a connection to the database
func Connect(dsn string) (*DB, error) {
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully connected to database")

	return &DB{sqlDB}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

