package postgresql

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	// "github.com/google/uuid"
	// "time"
)

func (p *PostgreSQL) DeleteUser(ctx context.Context, user *models.User) error {
	const op = "DeleteUser.DeleteUser"

	query := `DELETE FROM users WHERE id = $1`
	_, err := p.Db.ExecContext(ctx, query, user.ID)

	if err != nil {
		return fmt.Errorf("op: %v, error: %v", op, err)
	}

	return nil
}

func (p *PostgreSQL) SoftDeleteUser(ctx context.Context, user *models.User) error {
	const op = "SoftDeleteUser.SoftDeleteUser"
	query := `UPDATE users SET is_deleted = true WHERE id = $1`
	res, err := p.Db.ExecContext(ctx, query, user.ID.String())
	if err != nil {
		return fmt.Errorf("op: %v, error: %v", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("op: %v, error: %v", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("op: %v, rows: %v", op, user.ID)
	}

	return nil
}
