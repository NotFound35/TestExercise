package postgresql

import (
	"awesomeProject/config"
	"awesomeProject/internal/domain/models"
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

func (p *PostgreSQL) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		age INTEGER NOT NULL CHECK (age > 0),
		recording_date BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := p.db.Exec(query)
	if err != nil {
		p.logger.Error("ОШИБКА создания таблицы", zap.Error(err))
		return err
	}

	p.logger.Info("таблица СОЗДАНА!!!")
	return nil
}

func (p *PostgreSQL) SaveUser(user *models.User) error {
	query := `INSERT INTO users (id, first_name, last_name, age, recording_date) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := p.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Age,
		user.RecordingDate,
	)
	if err != nil {
		p.logger.Error("ОШИБКА сохранения юзера", zap.Error(err))
		return err
	}

	p.logger.Info("юзер сохранен - УСПЕХ!!!")
	return nil
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
