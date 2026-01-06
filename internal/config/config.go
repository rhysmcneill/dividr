package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	LogLevel    string
}

func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		slog.Warn(".env file not found, relying on environment variables")
	}

	return &Config{
		AppEnv:      getEnv("APP_ENV", "dev"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: buildDatabaseURL(),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// REDACTED LogValue hides sensitive fields when logging the config
func (c *Config) LogValue() slog.Value {
	// Return a group that hides sensitive fields
	return slog.GroupValue(
		slog.String("AppEnv", c.AppEnv),
		slog.String("Port", c.Port),
		slog.String("LogLevel", c.LogLevel),
		slog.String("DatabaseURL", "[REDACTED]"), // Hide the URL
	)
}

// buildDatabaseURL constructs the DB URL from env vars or Docker secret
func buildDatabaseURL() string {
	// If DATABASE_URL is already set, use it
	if dbURL := getEnv("DATABASE_URL", ""); dbURL != "" {
		return dbURL
	}

	// Otherwise, build from components
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "dividr")
	dbName := getEnv("DB_NAME", "dividr")
	password := getPassword()

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName)
}

// getPassword reads from Docker secret file or falls back to env var
func getPassword() string {
	// Try Docker secret first
	secretPath := "/run/secrets/db_password"
	if data, err := os.ReadFile(secretPath); err == nil {
		slog.Info("database password loaded from Docker secret")
		return strings.TrimSpace(string(data))
	}

	// Fallback to env var
	slog.Info("database password loaded from environment variable")
	return getEnv("DB_PASSWORD", "dividr")
}
