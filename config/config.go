package config

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	Env             string
	UserServiceURL  string
	JWTSecret       string
}

// LoadConfig читает env переменные и возвращает struct Config
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("SERVER_PORT", "8080"),
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		Env:             getEnv("ENV", "development"),
		UserServiceURL:  getEnv("USER_SERVICE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", ""),
	}
}

// Helpers
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
