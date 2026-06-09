package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	AppPort    string
	AppBaseURL string

	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	DBSSLMode         string
	DBMaxConnections  int
	DBMaxIdleConns    int

	JWTSecret    string
	JWTExpiryHours int

	DeepSeekAPIKey      string
	DeepSeekBaseURL     string
	DeepSeekModel       string
	DeepSeekMaxTokens   int
	DeepSeekTimeoutSecs int

	FonnteAPIKey  string
	FonnteBaseURL string
	FonnteDevice  string
}

func Load() (*Config, error) {
	// Only load .env in development; ignore error if file missing
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),

		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBName:           getEnv("DB_NAME", "kasiraiai"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBSSLMode:        getEnv("DB_SSL_MODE", "disable"),
		DBMaxConnections: getEnvInt("DB_MAX_CONNECTIONS", 25),
		DBMaxIdleConns:   getEnvInt("DB_MAX_IDLE_CONNECTIONS", 5),

		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 72),

		DeepSeekAPIKey:      getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekBaseURL:     getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
		DeepSeekModel:       getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
		DeepSeekMaxTokens:   getEnvInt("DEEPSEEK_MAX_TOKENS", 1000),
		DeepSeekTimeoutSecs: getEnvInt("DEEPSEEK_TIMEOUT_SECONDS", 30),

		FonnteAPIKey:  getEnv("FONNTE_API_KEY", ""),
		FonnteBaseURL: getEnv("FONNTE_BASE_URL", "https://api.fonnte.com"),
		FonnteDevice:  getEnv("FONNTE_DEVICE", ""),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
