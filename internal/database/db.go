package database

import (
	"fmt"
	"log"

	"pets_rest/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Interface defines the interface for database operations
type Interface interface {
	Health() error
	Close() error
}

// DB holds the database connection
type DB struct {
	*sqlx.DB
}

// Connect creates a new database connection and runs migrations
func Connect(cfg *config.Config) (*DB, error) {
	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Health checks the database connection
func (db *DB) Health() error {
	return db.Ping()
}
