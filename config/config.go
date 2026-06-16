package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDriver   string
	ServerPort string
	ServerHost string
	LogLevel   string
}

// LoadConfig reads environment variables and returns Config struct
func LoadConfig() (*Config, error) {
	// Load .env file (optional, won't fail if not found)
	godotenv.Load()

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_api_db"),
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}

	return config, nil
}

// Helper function to get env var with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	if c.DBDriver == "postgres" {
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
		)
	}
	// MySQL format
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

