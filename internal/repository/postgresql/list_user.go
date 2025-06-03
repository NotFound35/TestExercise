package postgresql

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
)

func (p *PostgreSQL) ListUsersPostgreSQL(
	ctx context.Context,
	minAge, maxAge *int,
	startDate, endDate *int64,
) ([]models.User, error) {
	const op = "ListUsersPostgreSQL"
	query := `
		SELECT id, first_name, last_name, age, recording_date 
		FROM users 
		WHERE 
			($1::int IS NULL OR age >= $1) AND
			($2::int IS NULL OR age <= $2) AND
			($3::bigint IS NULL OR recording_date >= $3) AND
			($4::bigint IS NULL OR recording_date <= $4)
	`

	rows, err := p.Db.QueryContext(ctx, query, minAge, maxAge, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("op: %s, %w", op, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.RecordingDate,
		)
		if err != nil {
			return nil, fmt.Errorf("op: %s, %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}
