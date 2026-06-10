package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName               string
	AppEnv                string
	Port                  string
	SQLitePath            string
	CorsAllowOrigin       string
	RequestTimeoutSeconds int
	DefaultAdminUsername  string
	DefaultAdminPassword  string
	DefaultModelName      string
	DefaultModelProvider  string
	DefaultModelGroup     string
	DeepSeekBaseURL       string
	DeepSeekAPIKey        string
	DeepSeekChannelName   string
}

func Load() Config {
	if hasEnvFile(".env") || hasEnvFile("../.env") {
		if err := godotenv.Load(".env", "../.env"); err != nil {
			log.Printf("failed to load env file, using system environment only: %v", err)
		}
	}

	return Config{
		AppName:               getEnv("APP_NAME", "AI Smart Chat Platform"),
		AppEnv:                getEnv("APP_ENV", "development"),
		Port:                  getEnv("BACKEND_PORT", "8080"),
		SQLitePath:            getEnv("SQLITE_PATH", "./data/app.db"),
		CorsAllowOrigin:       getEnv("CORS_ALLOW_ORIGIN", "http://localhost:5173"),
		RequestTimeoutSeconds: getEnvAsInt("REQUEST_TIMEOUT_SECONDS", 60),
		DefaultAdminUsername:  getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
		DefaultAdminPassword:  getEnv("DEFAULT_ADMIN_PASSWORD", "admin123456"),
		DefaultModelName:      getEnv("DEFAULT_MODEL_NAME", "deepseek-chat"),
		DefaultModelProvider:  getEnv("DEFAULT_MODEL_PROVIDER", "deepseek"),
		DefaultModelGroup:     getEnv("DEFAULT_MODEL_GROUP", "codex-group"),
		DeepSeekBaseURL:       getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
		DeepSeekAPIKey:        getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekChannelName:   getEnv("DEEPSEEK_CHANNEL_NAME", "DeepSeek Default Channel"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func hasEnvFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
