package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName              string
	AppEnv               string
	Port                 string
	SQLitePath           string
	CorsAllowOrigin      string
	DefaultAdminUsername string
	DefaultAdminPassword string
	DefaultModelName     string
	DefaultModelProvider string
	DefaultModelGroup    string
}

func Load() Config {
	if hasEnvFile(".env") || hasEnvFile("../.env") {
		if err := godotenv.Load(".env", "../.env"); err != nil {
			log.Printf("failed to load env file, using system environment only: %v", err)
		}
	}

	return Config{
		AppName:              getEnv("APP_NAME", "AI Smart Chat Platform"),
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 getEnv("BACKEND_PORT", "8080"),
		SQLitePath:           getEnv("SQLITE_PATH", "./data/app.db"),
		CorsAllowOrigin:      getEnv("CORS_ALLOW_ORIGIN", "http://localhost:5173"),
		DefaultAdminUsername: getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "admin123456"),
		DefaultModelName:     getEnv("DEFAULT_MODEL_NAME", "gpt-5.5"),
		DefaultModelProvider: getEnv("DEFAULT_MODEL_PROVIDER", "openai-compatible"),
		DefaultModelGroup:    getEnv("DEFAULT_MODEL_GROUP", "codex专用分组"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func hasEnvFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
