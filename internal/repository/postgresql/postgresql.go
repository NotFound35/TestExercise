package postgresql

import (
	"awesomeProject/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgreSQL struct {
	db     *sql.DB
	logger *zap.Logger
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

	return &PostgreSQL{db: db, logger: logger}, nil
}

func (p *PostgreSQL) Close() error {
	const op = "Close"
	if err := p.db.Close(); err != nil {
		fmt.Errorf("метод %v: %v", op, err)
		return err
	}
	p.logger.Info("соединения с БД закрыто - УСПЕХ!!!")
	return nil
}
