// закгрузка настроек подключения к PostgreSQL из переменных окружения .env
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
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
	configPath := "config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalln("ОШИБКА - файла конфига не существует")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalln("ОШИБКА - конфиг не читается")
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
	OutputPaths []string `yaml:"output_paths"` //сделать так: чтобы логи сохранялись в текстовый файл в корне проекта
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}
