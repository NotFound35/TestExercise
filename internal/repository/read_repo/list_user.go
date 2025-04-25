package read_repo

import (
	"awesomeProject/internal/domain/models"
	sq "github.com/Masterminds/squirrel"
)

func (r *ReadRepository) ListUsers(
	minAge, maxAge *int,
	minDate, maxDate *int64,
) ([]models.User, int, error) {
	qb := r.sb.
		Select("id", "first_name", "last_name", "age", "recording_date").
		From("users")

	if minAge != nil {
		qb = qb.Where(sq.GtOrEq{"age": *minAge})
	}
	if maxAge != nil {
		qb = qb.Where(sq.LtOrEq{"age": *maxAge})
	}
	if minDate != nil {
		qb = qb.Where(sq.GtOrEq{"recording_date": *minDate})
	}
	if maxDate != nil {
		qb = qb.Where(sq.LtOrEq{"recording_date": *maxDate})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.RecordingDate,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	// Получаем общее количество
	countQuery, _, err := r.sb.Select("COUNT(*)").From("users").ToSql()
	var total int
	err = r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
