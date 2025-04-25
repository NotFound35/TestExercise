package postgresql

import (
	"go.uber.org/zap"
)

func (p *PostgreSQL) CreateTables() error {
	p.logger.Info("начало миграций")
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

func Migrate(db *PostgreSQL) func() error {
	return db.CreateTables
}
