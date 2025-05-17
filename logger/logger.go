// Нужен, что бы аккуратно записывать, что происходит в программе (как черный ящик в самолете)
package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type LoggerConfig struct {
	Level       string   `yaml:"level"`
	OutputPaths []string `yaml:"output_paths"`
}

// InitLogger инициализирует zap.Logger на основе конфига
func InitLogger(cfg *LoggerConfig) (*zap.Logger, error) {
	const op = "InitLogger"

	// Настраиваем уровень логирования
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, err
	}

	// создание конфигурации для текстового вывода
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// создание файлового вывода
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("функция %s: %w", op, err)
	}

	// настройка вывода в файл и консоль
	fileWriteSyncer := zapcore.AddSync(file)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			fileWriteSyncer,
			level,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriteSyncer,
			level,
		),
	)

	return zap.New(core, zap.AddCaller()), nil
}

func New() *zap.Logger {
	// создание простого логгера с выводом в файл и консоль
	cfg := &LoggerConfig{
		Level:       "info",
		OutputPaths: []string{"logs.txt", "stdout"},
	}
	logger, _ := InitLogger(cfg)
	return logger
}
