package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Env  string
	Port string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// Storage configuration
	// Storage Configuration
	StorageProvider    string
	StorageEndpoint    string
	StorageRegion      string
	StorageBucket      string
	StorageAccessKey   string
	StorageSecretKey   string
	StorageUseSSL      bool
	StorageCredentials string

	// MinIO / S3
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL    bool
	MinIOBucket    string

	// JWT
	JWTSecret         string
	JWTExpirationTime time.Duration

	// Email
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string

	// URLs
	BaseURL     string
	FrontendURL string

	// File Upload
	MaxFileSize      string
	AllowedFileTypes string

	// Google OAuth2
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Env:  getEnv("ENV", "development"),
		Port: getEnv("PORT", "8080"),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://pets_user:pets_password@localhost:5432/pets_search?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),

		// Storage configuration
		StorageProvider:    getEnv("STORAGE_PROVIDER", "minio"),
		StorageEndpoint:    getEnv("STORAGE_ENDPOINT", getEnv("MINIO_ENDPOINT", "localhost:9000")),
		StorageRegion:      getEnv("STORAGE_REGION", "us-east-1"),
		StorageBucket:      getEnv("STORAGE_BUCKET", getEnv("MINIO_BUCKET", "pets-photos")),
		StorageAccessKey:   getEnv("STORAGE_ACCESS_KEY", getEnv("MINIO_ACCESS_KEY", "minioadmin")),
		StorageSecretKey:   getEnv("STORAGE_SECRET_KEY", getEnv("MINIO_SECRET_KEY", "minioadmin")),
		StorageUseSSL:      getEnvAsBool("STORAGE_USE_SSL", getEnvAsBool("MINIO_USE_SSL", false)),
		StorageCredentials: getEnv("STORAGE_CREDENTIALS", ""),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOUseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		MinIOBucket:    getEnv("MINIO_BUCKET", "pets-photos"),

		JWTSecret:         getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
		JWTExpirationTime: time.Duration(getEnvAsInt("JWT_EXPIRATION_TIME", 86400)) * time.Second,

		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@yourpetsearch.com"),

		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		MaxFileSize:        getEnv("MAX_FILE_SIZE", "10MB"),
		AllowedFileTypes:   getEnv("ALLOWED_FILE_TYPES", "jpg,jpeg,png,gif,pdf"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
