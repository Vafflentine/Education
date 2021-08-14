package apiserver

import (
	"Education/internal/app/db"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ServerConfig struct {
	Address  string
	Port     string
	DBConfig *db.DBConfig
}

// New returns a new Config struct
func NewConfig() *ServerConfig {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	return &ServerConfig{
		Address:  getEnv("ADDRESS", "127.0.0.1"),
		Port:     getEnv("SERVER_PORT", ""),
		DBConfig: db.NewConfig(),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
