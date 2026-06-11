package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config centralizes runtime settings read from .env or system environment.
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

// Load reads optional .env files first, then falls back to process environment
// and safe local-development defaults.
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

// getEnv returns the environment value or a fallback when the key is empty.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// getEnvAsInt parses an integer environment value, keeping the fallback on errors.
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

// hasEnvFile checks whether a candidate .env path exists before loading it.
func hasEnvFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
