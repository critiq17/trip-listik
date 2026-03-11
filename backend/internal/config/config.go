package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN                   string
	APIKeyAI              string
	HTTPAddr              string
	Env                   string
	JWTSecret             string
	TelegramWebAppSecret  string
	SupabaseURL           string
	SupabaseAnonKey       string
	SupabaseServiceKey    string
	SupabaseJWTSecret     string
	SupabaseStorageBucket string
	CORSOrigins           string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("File .env not found")
	}

	return &Config{
		DSN:                   os.Getenv("DSN"),
		APIKeyAI:              os.Getenv("API_KEY_AI"),
		HTTPAddr:              getHTTPAddr(),
		Env:                   getEnv("APP_ENV", "dev"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		TelegramWebAppSecret:  os.Getenv("TELEGRAM_WEBAPP_SECRET"),
		SupabaseURL:           os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:       os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceKey:    os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SupabaseJWTSecret:     os.Getenv("SUPABASE_JWT_SECRET"),
		SupabaseStorageBucket: getEnv("SUPABASE_STORAGE_BUCKET", "trip-photos"),
		CORSOrigins:           getEnv("CORS_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getHTTPAddr() string {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		return v
	}
	if v := os.Getenv("PORT"); v != "" {
		return ":" + v
	}
	return ":8080"
}

func (c *Config) ValidateAPI() error {
	if c.DSN == "" {
		return fmt.Errorf("DSN is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.TelegramWebAppSecret == "" {
		return fmt.Errorf("TELEGRAM_WEBAPP_SECRET is required")
	}
	return nil
}
