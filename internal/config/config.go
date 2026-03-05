package config

import "github.com/joho/godotenv"

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddr      string
	RedisPassword  string
	ServerPort     string
	WorkerCount    int
	MaxRetries     int
	RequestTimeout int
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6380"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		WorkerCount:    getEnvInt("WORKER_COUNT", 1),
		MaxRetries:     getEnvInt("MAX_RETRIES", 3),
		RequestTimeout: getEnvInt("REQUEST_TIMEOUT", 10),
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}
