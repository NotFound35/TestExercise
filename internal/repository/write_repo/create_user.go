package write_repo

import (
	"awesomeProject/internal/domain/models"
	"database/sql"
	"go.uber.org/zap"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository - конструктор
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// SaveUser - сохранение пользователя
func (r *PostgreSQL) SaveUser(user *models.User) error {
	query := `INSERT INTO users (id, first_name, last_name, age, recording_date) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Age,
		user.RecordingDate,
	)
	return err
}

type PostgreSQL struct {
	db     *sql.DB
	logger *zap.Logger
}
