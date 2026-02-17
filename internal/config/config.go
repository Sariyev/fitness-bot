package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	Debug         bool
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	AdminPort     string
	AdminAPIKey   string
	WebAppPort    string
	WebAppURL     string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		Debug:         getEnv("DEBUG", "false") == "true",
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "fitbot"),
		DBPassword:    getEnv("DB_PASSWORD", "fitbot_password"),
		DBName:        getEnv("DB_NAME", "fitness_bot"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		AdminPort:     getEnv("ADMIN_PORT", "8080"),
		AdminAPIKey:   getEnv("ADMIN_API_KEY", ""),
		WebAppPort:    getEnv("WEBAPP_PORT", "8081"),
		WebAppURL:     getEnv("WEBAPP_URL", ""),
	}
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
