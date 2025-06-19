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
