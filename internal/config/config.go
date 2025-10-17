package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	// Server config
	ServerHost         string        `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	ServerPort         int           `env:"SERVER_PORT" envDefault:"8080"`
	ServerReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"30s"`
	ServerWriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"30s"`

	// Database config
	DBHost            string        `env:"DB_HOST" envDefault:"localhost"`
	DBPort            int           `env:"DB_PORT" envDefault:"5432"`
	DBUser            string        `env:"DB_USER" envDefault:"postgres"`
	DBPassword        string        `env:"DB_PASSWORD" envDefault:""`
	DBName            string        `env:"DB_NAME" envDefault:"booking_db"`
	DBSSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`

	// JWT config
	JWTSecret string `env:"JWT_SECRET" envDefault:""`

	// Log config
	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"json"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	// Debug: Print all environment variables
	fmt.Println("=== DEBUG: Environment Variables ===")
	fmt.Printf("DB_PASSWORD: %s\n", os.Getenv("DB_PASSWORD"))
	fmt.Printf("JWT_SECRET: %s\n", os.Getenv("JWT_SECRET"))
	fmt.Println("=== All env vars ===")
	for _, env := range os.Environ() {
		if len(env) > 50 { // Only show first 50 chars to avoid secrets in logs
			fmt.Printf("%.50s...\n", env)
		} else {
			fmt.Println(env)
		}
	}
	fmt.Println("=== End DEBUG ===")

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Debug: Print parsed config
	fmt.Println("=== DEBUG: Parsed Config ===")
	fmt.Printf("cfg.DBPassword: '%s'\n", cfg.DBPassword)
	fmt.Printf("cfg.JWTSecret: '%s'\n", cfg.JWTSecret)
	fmt.Println("=== End Parsed Config ===")

	// Validate required fields
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
