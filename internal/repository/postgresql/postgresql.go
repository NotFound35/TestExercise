package postgresql

import (
	"awesomeProject/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgreSQL struct {
	Db     *sql.DB
	Logger *zap.Logger
}

func NewPostgreSQL(cfg *config.Config, logger *zap.Logger) (*PostgreSQL, error) {
	const op = "NewPostgreSQL"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("функция %v: %v", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("функция %v: %v ", op, err)
	}

	logger.Info("успешный коннект с PostgreSQL")

	return &PostgreSQL{Db: db, Logger: logger}, nil
}

func (p *PostgreSQL) Close() error {
	const op = "Close"
	if err := p.Db.Close(); err != nil {
		//todo log
		fmt.Errorf("метод %v: %v", op, err)
		return err
	}
	p.Logger.Info("соединения с БД закрыто - УСПЕХ!!!")
	return nil
}
