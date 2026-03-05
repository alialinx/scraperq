package config

import (
	"github.com/joho/godotenv"
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
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		RedisAddr:      getEnv("REDIS_ADDR", ""),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		WorkerCount:    getEnvInt("WORKER_COUNT", 3),
		MaxRetries:     getEnvInt("MAX_RETRIES", 3),
		RequestTimeout: getEnvInt("REQUEST_TIMEOUT", 10),
		DBHost:         getEnv("POSTGRES_HOST", ""),
		DBPort:         getEnv("POSTGRES_PORT", ""),
		DBUser:         getEnv("POSTGRES_USER", ""),
		DBPassword:     getEnv("POSTGRES_PASSWORD", ""),
		DBName:         getEnv("POSTGRES_DB", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
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
