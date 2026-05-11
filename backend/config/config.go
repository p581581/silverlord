package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	AdminPass  string
	Port       string
	// SMTP
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPassword string
	SmtpFrom     string
	SmtpFromName string
	SiteURL      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env not found, using system environment variables")
	}

	return &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "silverlord"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me-in-production"),
		AdminPass:    getEnv("ADMIN_PASSWORD", "admin123"),
		Port:         getEnv("PORT", "8080"),
		SmtpHost:     getEnv("SMTP_HOST", ""),
		SmtpPort:     getEnvInt("SMTP_PORT", 587),
		SmtpUser:     getEnv("SMTP_USER", ""),
		SmtpPassword: getEnv("SMTP_PASSWORD", ""),
		SmtpFrom:     getEnv("SMTP_FROM", ""),
		SmtpFromName: getEnv("SMTP_FROM_NAME", ""),
		SiteURL:      getEnv("SITE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
