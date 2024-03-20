package config

import (
	"os"
)

type Config struct {
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
}

func NewConfig() *Config {
	return &Config {
		DBHost: getEnv("DBHost", "localhost"),
		DBPort: getEnv("DBPort", "5432"),
		DBUser: getEnv("DBUser", "postgres"),
		DBPassword: getEnv("DBPassword", "postgres"),
		DBName: getEnv("DBName", "postgres"),
	}
}

func getEnv(key, defaultV string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultV
	}
	return value
}
