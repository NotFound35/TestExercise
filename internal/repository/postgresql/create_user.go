package postgresql

import (
	"awesomeProject/internal/domain/models"
	"fmt"
)

func (p *PostgreSQL) SaveUser(user *models.User) error {
	const op = "CreateUser.SaveUser"
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
		fmt.Errorf("метод %v: %v", op, err)
		return err
	}

	p.logger.Info("юзер сохранен - УСПЕХ!!!")
	return nil
}
