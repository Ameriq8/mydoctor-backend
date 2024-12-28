package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Driver   string
}

type LoadedConfig struct {
	Config
	DatabaseConfig
}

func LoadConfig() LoadedConfig {
	return LoadedConfig{
		Config: Config{
			ServerPort: getEnv("SERVER_PORT", "8080"),
		},
		DatabaseConfig: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3002),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "mydb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Driver:   getEnv("DB_DRIVER", "postgres"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
