// закгрузка настроек подключения к PostgreSQL из переменных окружения .env
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	Server ServerConfig
}

type ServerConfig struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     port,
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic("НЕ ПОЛУЧИЛАСЬ загрузка конфига")
	}
	return cfg
}
