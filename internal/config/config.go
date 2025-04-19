package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

// Server holds server configuration
type Server struct {
	Address string
}

// Database holds database configuration
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWT holds JWT configuration
type JWT struct {
	Secret           string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

// OAuth holds OAuth provider configuration
type OAuth struct {
	// Google OAuth configuration
	GoogleClientID     string
	GoogleClientSecret string
	// Facebook OAuth configuration
	FacebookAppID     string
	FacebookAppSecret string
}

// Logging holds logging configuration
type Logging struct {
	Level zapcore.Level
}

// Config holds all application configuration
type Config struct {
	Server   Server
	Database Database
	JWT      JWT
	OAuth    OAuth
	Logging  Logging
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Server: Server{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "myapp"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWT{
			Secret:           getEnv("JWT_SECRET", "your-secret-key"),
			AccessExpiresIn:  time.Duration(getEnvAsInt("JWT_ACCESS_EXPIRES_IN", 24)) * time.Hour,
			RefreshExpiresIn: time.Duration(getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 7)) * 24 * time.Hour,
		},
		OAuth: OAuth{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			FacebookAppID:      getEnv("FACEBOOK_APP_ID", ""),
			FacebookAppSecret:  getEnv("FACEBOOK_APP_SECRET", ""),
		},
		Logging: Logging{
			Level: getLogLevel(getEnv("LOG_LEVEL", "info")),
		},
	}, nil
}

// getEnv retrieves environment variables with fallback values
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves environment variables as integers with fallback values
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return int(value.Hours())
	}
	return fallback
}

// getLogLevel converts a string log level to zapcore.Level
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
