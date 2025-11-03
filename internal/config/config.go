package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	ServerPort string
	DataBase   DataBaseConfig
}

// DataBaseConfig holds the database configuration
type DataBaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLMode  string
	MaxConn  int32
	MinConn  int32
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Parse connection limits with defaults
	maxConn, err := strconv.ParseInt(os.Getenv("DB_MAX_CONNS"), 10, 32)
	if err != nil {
		maxConn = 5 // Default value
	}

	minConn, err := strconv.ParseInt(os.Getenv("DB_MIN_CONNS"), 10, 32)
	if err != nil {
		minConn = 1 // Default value
	}

	config := &Config{
		ServerPort: getEnvWithDefault("PORT", "8080"),
		DataBase: DataBaseConfig{
			Host:     getEnvWithDefault("DB_HOST", "localhost"),
			Port:     getEnvWithDefault("DB_PORT", "5432"),
			User:     getEnvWithDefault("DB_USER", "postgres"),
			Password: getEnvWithDefault("DB_PASSWORD", "password"),
			DBname:   getEnvWithDefault("DB_NAME", "testdb"),
			SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
			MaxConn:  int32(maxConn),
			MinConn:  int32(minConn),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return fmt.Errorf("server port is required")
	}
	return c.DataBase.Validate()
}

// Validate checks if database configuration is valid
func (c *DataBaseConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if c.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.DBname == "" {
		return fmt.Errorf("database name is required")
	}
	if c.MaxConn < c.MinConn {
		return fmt.Errorf("max connections cannot be less than min connections")
	}
	return nil
}

func (c *DataBaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBname, c.SSLMode,
	)
}
