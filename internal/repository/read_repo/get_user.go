package read_repo

import (
	"awesomeProject/internal/domain/models"
	sq "github.com/Masterminds/squirrel"
)

func (r *ReadRepository) GetUser(id string) (*models.User, error) {
	query, args, err := r.sb.
		Select("id", "first_name", "last_name", "age", "recording_date").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := &models.User{}
	err = r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.RecordingDate,
	)
	return user, err
}
