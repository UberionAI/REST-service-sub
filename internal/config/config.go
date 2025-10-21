package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort      string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresHost string
	PostgresPort string
	PgSSLMode    string
}

//LoadConfig loads the config from the environment

func LoadConfig() *Config {
	cfg := &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		PostgresUser: getEnv("POSTGRES_USER", "postgres"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:   getEnv("POSTGRES_DB", "subscription_db"),
		PostgresHost: getEnv("POSTGRES_HOST", "127.0.0.1"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PgSSLMode:    getEnv("POSTGRES_SSLMODE", "disable"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.PostgresUser, c.PostgresPass, c.PostgresHost, c.PostgresPort, c.PostgresDB, c.PgSSLMode)
}
