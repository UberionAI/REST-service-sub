package config

import "github.com/joho/godotenv"

type Config struct {
	Port         string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresHost string
	PostgresPort string
	PgSSLMode    string
}

//LoadConfig loads the config from the environment

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
}
