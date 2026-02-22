package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

var AppConfig *Config

// LoadConfig loads environment variables and initializes the configuration
func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		Port:        getEnv("PORT", "8080"),
	}

	// Validate required environment variables
	if AppConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if AppConfig.JWTSecret == "default-secret-key" {
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable for production")
	}
}

// getEnv retrieves an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
