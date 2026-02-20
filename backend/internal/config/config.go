package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env string // development | staging | production
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnIdleTime time.Duration
}

type JWTConfig struct {
	Secret        string
	ExpiryMinutes int
}

type CORSConfig struct {
	AllowedOrigins []string
}

// Load reads configuration from environment variables (and optionally a .env file).
// Viper automatically reads env vars — no external file dependency in production.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("DB_MAX_CONNS", 20)
	v.SetDefault("DB_MIN_CONNS", 2)
	v.SetDefault("DB_MAX_CONN_IDLE_TIME", "30m")
	v.SetDefault("JWT_EXPIRY_MINUTES", 60)
	v.SetDefault("ALLOWED_ORIGINS", "http://localhost:5173")

	// Bind all environment variables automatically
	v.AutomaticEnv()

	// In development, also try loading from .env file (optional — does not fail if absent)
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig() // intentionally ignore error — .env is optional in production

	idleTime, err := time.ParseDuration(v.GetString("DB_MAX_CONN_IDLE_TIME"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONN_IDLE_TIME: %w", err)
	}

	origins := strings.Split(v.GetString("ALLOWED_ORIGINS"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	cfg := &Config{
		App: AppConfig{
			Env: v.GetString("APP_ENV"),
		},
		Server: ServerConfig{
			Host: v.GetString("SERVER_HOST"),
			Port: v.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			URL:             v.GetString("DATABASE_URL"),
			MaxConns:        int32(v.GetInt("DB_MAX_CONNS")),
			MinConns:        int32(v.GetInt("DB_MIN_CONNS")),
			MaxConnIdleTime: idleTime,
		},
		JWT: JWTConfig{
			Secret:        v.GetString("JWT_SECRET"),
			ExpiryMinutes: v.GetInt("JWT_EXPIRY_MINUTES"),
		},
		CORS: CORSConfig{
			AllowedOrigins: origins,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate ensures all required fields are present.
func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	if c.JWT.ExpiryMinutes < 15 || c.JWT.ExpiryMinutes > 60 {
		return fmt.Errorf("JWT_EXPIRY_MINUTES must be between 15 and 60")
	}
	return nil
}

// IsProd returns true when running in production mode.
func (c *Config) IsProd() bool {
	return c.App.Env == "production"
}
