package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper loads configuration with flexible source priority:
// 1. Environment variables (highest priority)
// 2. .env file (if exists)
// 3. config.json file (fallback)
//
// This allows users to choose their preferred configuration method:
// - Use .env for Docker/containerized environments
// - Use config.json for traditional JSON configuration
// - Mix both: .env overrides config.json values
func NewViper() *viper.Viper {
	config := viper.New()

	// Set default values
	setDefaults(config)

	// Try to load from .env file first (optional)
	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	config.AddConfigPath("./../")

	// AutomaticEnv allows environment variables to override config files
	config.AutomaticEnv()

	// Try reading .env file (ignore error if not exists)
	if err := config.ReadInConfig(); err != nil {
		// If .env not found, try config.json
		config.SetConfigName("config")
		config.SetConfigType("json")
		config.AddConfigPath(".")
		config.AddConfigPath("./../")

		if err := config.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error: no config file found (.env or config.json): %w", err))
		}
	}

	// Map environment variable names to config keys
	// This allows both DB_HOST and database.host to work
	config.BindEnv("app.name", "APP_NAME")
	config.BindEnv("app.env", "APP_ENV")
	config.BindEnv("web.port", "WEB_PORT")
	config.BindEnv("web.prefork", "WEB_PREFORK")
	config.BindEnv("api.prefix", "API_PREFIX")
	config.BindEnv("api.version", "API_VERSION")
	config.BindEnv("database.host", "DB_HOST")
	config.BindEnv("database.port", "DB_PORT")
	config.BindEnv("database.username", "DB_USERNAME")
	config.BindEnv("database.password", "DB_PASSWORD")
	config.BindEnv("database.name", "DB_NAME")
	config.BindEnv("database.pool.idle", "DB_POOL_IDLE")
	config.BindEnv("database.pool.max", "DB_POOL_MAX")
	config.BindEnv("database.pool.lifetime", "DB_POOL_LIFETIME")
	config.BindEnv("jwt.secret", "JWT_SECRET")
	config.BindEnv("jwt.access_expire_minutes", "JWT_ACCESS_EXPIRE_MINUTES")
	config.BindEnv("jwt.refresh_expire_days", "JWT_REFRESH_EXPIRE_DAYS")
	config.BindEnv("cors.allowed_origins", "CORS_ALLOWED_ORIGINS")
	config.BindEnv("cors.allowed_methods", "CORS_ALLOWED_METHODS")
	config.BindEnv("cors.allowed_headers", "CORS_ALLOWED_HEADERS")
	config.BindEnv("rate_limit.enabled", "RATE_LIMIT_ENABLED")
	config.BindEnv("rate_limit.rpm", "RATE_LIMIT_RPM")
	config.BindEnv("log.level", "LOG_LEVEL")
	config.BindEnv("email.host", "EMAIL_HOST")
	config.BindEnv("email.port", "EMAIL_PORT")
	config.BindEnv("email.username", "EMAIL_USERNAME")
	config.BindEnv("email.password", "EMAIL_PASSWORD")
	config.BindEnv("email.from", "EMAIL_FROM")
	config.BindEnv("base_url", "BASE_URL")

	return config
}

// setDefaults sets default configuration values
func setDefaults(config *viper.Viper) {
	// App defaults
	config.SetDefault("app.name", "go-clean-arch-saas")
	config.SetDefault("app.env", "development")

	// Web defaults
	config.SetDefault("web.port", 3000)
	config.SetDefault("web.prefork", false)

	// API defaults
	config.SetDefault("api.prefix", "/api")
	config.SetDefault("api.version", "v1")

	// Database defaults
	config.SetDefault("database.host", "localhost")
	config.SetDefault("database.port", 3306)
	config.SetDefault("database.username", "root")
	config.SetDefault("database.password", "")
	config.SetDefault("database.name", "go_clean_arch_saas")
	config.SetDefault("database.pool.idle", 10)
	config.SetDefault("database.pool.max", 100)
	config.SetDefault("database.pool.lifetime", 300)

	// JWT defaults
	config.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	config.SetDefault("jwt.access_expire_minutes", 60)
	config.SetDefault("jwt.refresh_expire_days", 7)

	// CORS defaults
	config.SetDefault("cors.allowed_origins", "http://localhost:3000,http://localhost:8080")
	config.SetDefault("cors.allowed_methods", "GET,POST,PUT,PATCH,DELETE")
	config.SetDefault("cors.allowed_headers", "Origin,Content-Type,Accept,Authorization")

	// Rate limit defaults
	config.SetDefault("rate_limit.enabled", false)
	config.SetDefault("rate_limit.rpm", 1000)

	// Logging defaults
	config.SetDefault("log.level", 6)

	// Email defaults (empty means email disabled in development)
	config.SetDefault("email.host", "")
	config.SetDefault("email.port", 587)
	config.SetDefault("email.username", "")
	config.SetDefault("email.password", "")
	config.SetDefault("email.from", "noreply@localhost")

	// Base URL default
	config.SetDefault("base_url", "http://localhost:3000")
}
