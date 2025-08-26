package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            int
	GinMode         string

	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string

	MCPXHSEndpoint  string
	MCPMapsEndpoint string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func Load() Config {
	return Config{
		Port:            getEnvInt("PORT", 8080),
		GinMode:         getEnv("GIN_MODE", "release"),

		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnvInt("DB_PORT", 3306),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", "root"),
		DBName:          getEnv("DB_NAME", "photodb"),
		DBSSLMode:       getEnv("DB_SSLMODE", ""),

		MCPXHSEndpoint:  getEnv("MCP_XHS_ENDPOINT", "http://localhost:9001"),
		MCPMapsEndpoint: getEnv("MCP_MAPS_ENDPOINT", "http://localhost:9002"),
	}
} 