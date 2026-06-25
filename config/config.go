package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort  string
	QueueBuffer int
	WorkerCount int
	RedisURL    string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:  getEnv("AEGIS_PORT", "8080"),
		QueueBuffer: getEnvAsInt("AEGIS_QUEUE_BUFFER", 10000),
		WorkerCount: getEnvAsInt("AEGIS_WORKER_COUNT", 10),
		RedisURL:    getEnv("AEGIS_PORT", "redis://localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
