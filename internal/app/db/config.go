package db

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DBConfig struct {
	DBUser     string
	DBPassword string
	Address    string
	DBPort     string
	DBName     string
	Dns        string
}

// New returns a new Config struct
func NewConfig() *DBConfig {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	return &DBConfig{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		Address:    getEnv("ADDRESS", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
