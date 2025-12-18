package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is a wrapper around sqlx.DB
type DB struct {
	*sqlx.DB
}

// Connect creates a new database connection
func Connect(databaseURL string) (*DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Successfully connected to database")
	return &DB{db}, nil
}

// RunMigrations runs all migration files in the migrations directory
func (db *DB) RunMigrations(migrationsPath string) error {
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort files to ensure they run in order
	sort.Strings(files)

	for _, file := range files {
		log.Printf("Running migration: %s", filepath.Base(file))
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// BeginTx starts a new transaction
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.DB.Begin()
}
