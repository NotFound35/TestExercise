package postgresql

import (
	"awesomeProject/internal/domain/models"
)

func (r *PostgreSQL) ListUsers(
	minAge, maxAge *int,
	minDate, maxDate *int64,
) ([]models.User, int, error) {
	panic("implement me")
}
