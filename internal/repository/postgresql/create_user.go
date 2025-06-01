package postgresql

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (p *PostgreSQL) SaveUser(ctx context.Context, user *models.User) error {
	const op = "CreateUser.SaveUser"

	query := `INSERT INTO users (id, first_name, last_name, age, recording_date) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := p.Db.ExecContext(ctx, query,
		uuid.New(),
		user.FirstName,
		user.LastName,
		user.Age,
		time.Now().Unix(),
	)

	if err != nil {
		return fmt.Errorf("op: %v, error: %v", op, err)
	}

	p.Logger.Info("юзер сохранен - УСПЕХ!!!")
	return nil
}
