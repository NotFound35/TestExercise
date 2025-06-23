package postgresql

import (
	"fmt"
)

func (p *PostgreSQL) CreateTables() error {
	const op = "CreateTables"
	p.Logger.Info("начало миграций")
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		age INTEGER NOT NULL CHECK (age > 0),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    is_deleted BOOLEAN DEFAULT FALSE
	)`

	_, err := p.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("метод %v: %v", op, err)
	}

	p.Logger.Info("таблица СОЗДАНА!!!")
	return nil
}

func Migrate(db *PostgreSQL) {
	fmt.Println("migrations success")
	db.CreateTables()
}
