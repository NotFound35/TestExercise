package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env"`
	Database   Database   `yaml:"database"`
	Logger     zap.Logger `yaml:"logger"`
	HTTPServer HTTPServer `yaml:"http-server"`
}

func MustLoad() (*Config, error) {
	const op = "MustLoad"
	configPath := "config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("op %s: %w", op, err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("op %s: %w", op, err)
	}

	return &cfg, nil
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl_mode"`
}
type Logger struct {
	Level       string   `yaml:"level"`
	Encoding    string   `yaml:"encoding"`
	OutputPaths []string `yaml:"output_paths"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}
