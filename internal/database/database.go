package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/vention/booking_api/internal/config"
)

type DB struct {
	*sql.DB
}

func NewPostgreSQL(cfg *config.Config, logger zerolog.Logger) (*DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("Successfully connected to PostgreSQL database")

	return &DB{db}, nil
}

func (db *DB) Close() {
	db.DB.Close()
}

func (db *DB) Health() error {
	return db.Ping()
}
