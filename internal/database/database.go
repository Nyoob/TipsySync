package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the sql.DB connection and provides helper methods
type DB struct {
	*sql.DB
}

// New creates a new database instance
func NewDatabase() *DB {
	cfgDir, _ := os.UserConfigDir()
	dbDir := cfgDir + "/TipAggregator"
	dbPath := dbDir + "/aggregator.db"
	// dbDir := "./data"
	// dbPath := "./data/db.db"
	os.MkdirAll(dbDir, 0755)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("failed to open database: %w", err)
		return nil
	}

	// Test connection
	if err := db.Ping(); err != nil {
		fmt.Println("failed to ping database: %w", err)
		return nil
	}

	dbInstance := &DB{DB: db}

	// Run migrations
	if err := dbInstance.migrate(); err != nil {
		fmt.Println("failed to migrate database: %w", err)
		return nil
	}

	return dbInstance
}

// migrate runs database migrations
func (db *DB) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS settings (
            id TEXT PRIMARY KEY,
            value TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS settings_providers (
            id TEXT PRIMARY KEY,
						enabled BOOLEAN,
            apiToken TEXT,
						fetchInterval INTEGER
        );`,
		`CREATE TABLE IF NOT EXISTS event (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            type TEXT,
						data TEXT
        );`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
