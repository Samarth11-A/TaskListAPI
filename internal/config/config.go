package config

import (
	"os"
	"strconv"

	"github.com/Samarth11-A/TaskListAPI/internal/database"
)

// Config holds application configuration
type Config struct {
	ServerPort string
	DB         database.Config
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return Config{
		ServerPort: getEnv("SERVER_PORT", ":50051"),
		DB: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "tasklist"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
