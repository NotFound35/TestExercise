package postgresql

import (
	"fmt"
)

func (p *PostgreSQL) CreateTables() error {
	const op = "CreateTables"
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		age INTEGER NOT NULL CHECK (age > 0),
		recording_date BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    is_deleted BOOLEAN DEFAULT FALSE
	)`

	_, err := p.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("op %s: err %w", op, err)
	}

	return nil
}

func Migrate(db *PostgreSQL) {
	fmt.Println("migrations success")
	db.CreateTables()
}
