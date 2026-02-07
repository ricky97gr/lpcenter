package config

import (
	"os"
)

type Config struct {
	DBHost     string `json:"dbHost" yaml:"dbHost"`
	DBPort     string `json:"dbPort" yaml:"dbPort"`
	DBUser     string `json:"dbUser" yaml:"dbUser"`
	DBPassword string `json:"dbPassword" yaml:"dbPassword"`
	DBName     string `json:"dbName" yaml:"dbName"`
	ServerPort string `json:"serverPort" yaml:"serverPort"`
}

var c *Config

func GetConfig() (Config, error) {
	if c == nil {
		load()
	}
	if c == nil {
		return Config{}, nil
	}
	return *c, nil
}

func load() {
	c = &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "123456"),
		DBName:     getEnv("DB_NAME", "lpcenter"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
