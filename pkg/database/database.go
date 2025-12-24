// Package database provides database connection management and utilities using sqlx.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// DB wraps sqlx.DB with additional functionality
type DB struct {
	*sqlx.DB
	driver string
}

// Config holds database configuration
type Config struct {
	Driver          string        // postgres, mysql, sqlite
	DSN             string        // Data Source Name / Connection String
	MaxOpenConns    int           // Maximum number of open connections
	MaxIdleConns    int           // Maximum number of idle connections
	ConnMaxLifetime time.Duration // Maximum connection lifetime
	ConnMaxIdleTime time.Duration // Maximum connection idle time
}

// DefaultConfig returns a default database configuration
func DefaultConfig() Config {
	return Config{
		Driver:          "sqlite",
		DSN:             ":memory:",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// New creates a new database connection with the given configuration
func New(cfg Config) (*DB, error) {
	if cfg.Driver == "" {
		return nil, fmt.Errorf("database driver is required")
	}
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database DSN is required")
	}

	db, err := sqlx.Connect(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	return &DB{
		DB:     db,
		driver: cfg.Driver,
	}, nil
}

// HealthCheck performs a basic health check on the database connection
func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Driver returns the database driver name
func (db *DB) Driver() string {
	return db.driver
}

// WithTransaction executes a function within a database transaction
func (db *DB) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// RunMigrations executes database migrations
func (db *DB) RunMigrations(ctx context.Context, migrations []Migration) error {
	// Create migrations table if it doesn't exist
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	if db.driver == "postgres" {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migrations (
				version INTEGER PRIMARY KEY,
				description TEXT NOT NULL,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	} else if db.driver == "mysql" {
		createTableSQL = `
			CREATE TABLE IF NOT EXISTS migrations (
				version INTEGER PRIMARY KEY,
				description TEXT NOT NULL,
				applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)`
	}

	if _, err := db.ExecContext(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedVersions []int
	err := db.SelectContext(ctx, &appliedVersions, "SELECT version FROM migrations ORDER BY version")
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	appliedMap := make(map[int]bool)
	for _, v := range appliedVersions {
		appliedMap[v] = true
	}

	// Run pending migrations
	for _, migration := range migrations {
		if appliedMap[migration.Version] {
			continue
		}

		err := db.WithTransaction(ctx, func(tx *sqlx.Tx) error {
			// Execute migration
			if _, err := tx.ExecContext(ctx, migration.Up); err != nil {
				return fmt.Errorf("migration %d failed: %w", migration.Version, err)
			}

			// Record migration
			_, err := tx.ExecContext(ctx,
				"INSERT INTO migrations (version, description) VALUES (?, ?)",
				migration.Version, migration.Description)
			if err != nil && db.driver == "postgres" {
				_, err = tx.ExecContext(ctx,
					"INSERT INTO migrations (version, description) VALUES ($1, $2)",
					migration.Version, migration.Description)
			}

			if err != nil {
				return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
