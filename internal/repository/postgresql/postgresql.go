package postgresql

import (
	"awesomeProject/config"
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
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ОШИБКА открытия коннекта с БД")
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ОШИБКА соединения с БД")
	}

	logger.Info("успешный коннект с PostgreSQL")
	return &PostgreSQL{db: db, logger: logger}, nil
}

func (p *PostgreSQL) GetConnection() *sql.DB {
	return p.db
}

func (p *PostgreSQL) Close() error {
	if err := p.db.Close(); err != nil {
		p.logger.Error("ОШИБКА закрытия соединения с БД", zap.Error(err))
		return err
	}
	p.logger.Info("соединения с БД закрыто - УСПЕХ!!!")
	return nil
}
