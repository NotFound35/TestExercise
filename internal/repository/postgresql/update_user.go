package postgresql

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	// "github.com/google/uuid"
	// "time"
)

func (p *PostgreSQL) UserUpdate(ctx context.Context, user *models.User) error {
	const op = "UpdateUser"

	query := `
			UPDATE users
			SET first_name = $1, last_name = $2, age = $3
			WHERE id = $4 AND is_deleted = false`

	res, err := p.Db.ExecContext(ctx, query, user.FirstName, user.LastName, user.Age, user.ID)

	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Rows error %s %w", op, err)
	}

	if rows == 0 {
		return fmt.Errorf("User %s not found", user.FirstName)
	}

	return nil

}
