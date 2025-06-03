package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env"`
	Database   Database   `yaml:"database"`
	Logger     Logger     `yaml:"logger"`
	HTTPServer HTTPServer `yaml:"http-server"`
}

func MustLoad() *Config {
	const op = "MustLoad"
	configPath := "config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		//todo log / return
		fmt.Errorf("функция %v: %v", op, err)
		return nil
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		fmt.Errorf("функция %v: %v", op, err)
	}

	return &cfg
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
