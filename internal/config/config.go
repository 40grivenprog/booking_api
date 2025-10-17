package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Database DatabaseConfig `envPrefix:"DB_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Log      LogConfig      `envPrefix:"LOG_"`
}

type ServerConfig struct {
	Host         string        `env:"HOST" envDefault:"0.0.0.0"`
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"30s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
}

type DatabaseConfig struct {
	Host            string        `env:"HOST" envDefault:"localhost"`
	Port            int           `env:"PORT" envDefault:"5432"`
	User            string        `env:"USER" envDefault:"postgres"`
	Password        string        `env:"PASSWORD" envDefault:""`
	DBName          string        `env:"NAME" envDefault:"booking_db"`
	SSLMode         string        `env:"SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"5m"`
}

type JWTConfig struct {
	Secret string `env:"SECRET" envDefault:""`
}

type LogConfig struct {
	Level  string `env:"LEVEL" envDefault:"info"`
	Format string `env:"FORMAT" envDefault:"json"`
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

	// Validate required fields
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
