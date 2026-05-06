package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	Environment string
	JWTSecret   string
}

func Load() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// En Render, esto no debería pasar
		// Pero si pasa, usamos un default para desarrollo local
		databaseURL = "postgres://postgres:postgres@localhost:5432/todos_db"
	}
	return &Config{
		DatabaseURL: databaseURL,
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-super-secret-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
