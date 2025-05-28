package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	*sqlx.DB
}

type Config struct {
	Host     string
	Port     int
	Username string // Use one consistent field name
	Password string
	DBName   string // This should match what you use in the connection string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*PostgresDB, error) {
	// Fix connection string formatting
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}
	log.Printf("Connected to PostgreSQL database at %s:%d", cfg.Host, cfg.Port)
	return &PostgresDB{DB: db}, nil
}

func (db *PostgresDB) CreateTasksTable() error {
	// Make schema match your Task protobuf definition
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}
	log.Println("Tasks table created successfully")
	return nil
}

func (db *PostgresDB) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	log.Println("Database connection closed")
	return nil
}
