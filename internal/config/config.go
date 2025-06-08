package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/Samarth11-A/TaskListAPI/internal/database"
)

type AppConfig struct {
	Environment string
}

type ServerConfig struct {
	Port       string
	ServerName string
}

// Config holds application configuration
type Config struct {
	AppConfig AppConfig
	SConfig   ServerConfig
	DB        database.Config
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	// Load .env file if it exists
	loadEnvFile()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return Config{
		AppConfig: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		SConfig: ServerConfig{
			Port:       getEnv("SERVER_PORT", "50051"),
			ServerName: getEnv("SERVER_NAME", "localhost"),
		},
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

// loadEnvFile loads environment variables from .env file
func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		// .env file doesn't exist, continue with system env vars
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}

			// Only set if not already set in system environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
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
