package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	Port          string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

var AppConfig *Config

func LoadConfig() error {
	// Load .env file - thử nhiều đường dẫn
	envPaths := []string{
		".env",
		"./.env",
		filepath.Join(".", ".env"),
	}

	var envLoaded bool
	var lastErr error
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			log.Printf("✅ Loaded .env file from: %s", path)
			break
		} else {
			lastErr = err
		}
	}

	if !envLoaded {
		log.Printf("⚠️  Warning: .env file not found (last error: %v), using system environment variables", lastErr)
	}

	AppConfig = &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "ecommerce_db"),
		Port:          getEnv("PORT", "8080"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
	)
}
