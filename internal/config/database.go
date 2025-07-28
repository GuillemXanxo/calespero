package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDBConfigFromEnv creates a new database configuration from environment variables
func NewDBConfigFromEnv() *DBConfig {
	return &DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		DBName:   getEnvOrDefault("DB_NAME", "calespero"),
		SSLMode:  getEnvOrDefault("DB_SSL_MODE", "require"), // Use 'require' for Render, 'disable' for local development
	}
}

// ConnectionURL returns the PostgreSQL connection URL
func (c *DBConfig) ConnectionURL() string {
	// If DATABASE_URL is set (like in Render), use it directly
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// Otherwise, build the connection string from individual components
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

// Connect establishes a connection to the database
func (c *DBConfig) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging the database: %v", err)
	}

	return db, nil
}

// Helper function to get environment variable with a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
