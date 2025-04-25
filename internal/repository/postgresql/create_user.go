package postgresql

import (
	"awesomeProject/internal/domain/models"
	"go.uber.org/zap"
)

func (p *PostgreSQL) SaveUser(user *models.User) error {
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
		p.logger.Error("ОШИБКА сохранения юзера", zap.Error(err))
		return err
	}

	p.logger.Info("юзер сохранен - УСПЕХ!!!")
	return nil
}
