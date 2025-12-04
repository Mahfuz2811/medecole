package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Cleanup  CleanupConfig
	OAuth    OAuthConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int

	// Connection Pool Settings
	PoolSize      int           // Maximum number of socket connections (default: 10)
	MinIdleConns  int           // Minimum number of idle connections (default: 5)
	MaxRetries    int           // Maximum number of retries before giving up (default: 3)
	DialTimeout   time.Duration // Timeout for establishing new connections (default: 5s)
	ReadTimeout   time.Duration // Timeout for socket reads (default: 3s)
	WriteTimeout  time.Duration // Timeout for socket writes (default: 3s)
	IdleTimeout   time.Duration // Amount of time after which client closes idle connections (default: 5m)
	PoolTimeout   time.Duration // Amount of time client waits for connection if all are busy (default: 4s)
	MaxConnAge    time.Duration // Connection age at which client retires the connection (default: 0 = disabled)
	IdleCheckFreq time.Duration // Frequency of idle checks made by idle connections reaper (default: 1m)
}

// NewRedisConfigWithDefaults creates a RedisConfig with optimal default pool settings
func NewRedisConfigWithDefaults(host, port, password string, db int) RedisConfig {
	return RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,

		// Optimal default pool settings
		PoolSize:      10,              // Good for most applications
		MinIdleConns:  5,               // Keep some connections ready
		MaxRetries:    3,               // Reasonable retry count
		DialTimeout:   5 * time.Second, // Connection establishment timeout
		ReadTimeout:   3 * time.Second, // Socket read timeout
		WriteTimeout:  3 * time.Second, // Socket write timeout
		IdleTimeout:   5 * time.Minute, // Close idle connections after 5 minutes
		PoolTimeout:   4 * time.Second, // Wait for connection from pool
		MaxConnAge:    0,               // Disabled - let Redis handle connection lifecycle
		IdleCheckFreq: 1 * time.Minute, // Check for idle connections every minute
	}
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    string
	GinMode string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	FrontendURL string
}

// CleanupConfig holds background cleanup service configuration
type CleanupConfig struct {
	Enabled         bool          // Whether cleanup service is enabled (default: true)
	CleanupInterval time.Duration // How often to run cleanup (default: 1 minute)
	GracePeriod     time.Duration // Grace period to avoid race conditions (default: 2 minutes)
}

// OAuthConfig holds OAuth provider configurations
type OAuthConfig struct {
	Google   GoogleOAuthConfig
	Facebook FacebookOAuthConfig
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// FacebookOAuthConfig holds Facebook OAuth configuration
type FacebookOAuthConfig struct {
	AppID       string
	AppSecret   string
	RedirectURL string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Parse Redis DB number
	redisDB := 0
	if dbStr := getEnv("REDIS_DB", "0"); dbStr != "" {
		if parsed, err := strconv.Atoi(dbStr); err == nil {
			redisDB = parsed
		}
	}

	// Parse cleanup configuration
	cleanupEnabled := getEnv("CLEANUP_ENABLED", "true") == "true"
	cleanupInterval := parseDuration("CLEANUP_INTERVAL", "1m")
	gracePeriod := parseDuration("CLEANUP_GRACE_PERIOD", "2m")

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			Name:     getEnv("DB_NAME", "quizora"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		},
		CORS: CORSConfig{
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
		Cleanup: CleanupConfig{
			Enabled:         cleanupEnabled,
			CleanupInterval: cleanupInterval,
			GracePeriod:     gracePeriod,
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
			},
			Facebook: FacebookOAuthConfig{
				AppID:       getEnv("FACEBOOK_APP_ID", ""),
				AppSecret:   getEnv("FACEBOOK_APP_SECRET", ""),
				RedirectURL: getEnv("FACEBOOK_REDIRECT_URL", "http://localhost:8080/api/v1/auth/facebook/callback"),
			},
		},
	}
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses duration from environment variable with fallback to default
func parseDuration(key, defaultValue string) time.Duration {
	durationStr := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Printf("Invalid duration format for %s: %s, using default: %s", key, durationStr, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}
