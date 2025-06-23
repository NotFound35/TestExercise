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
		SELECT id, first_name, last_name, age, created_at 
		FROM users 
		WHERE 
		    is_deleted = false AND
			($1 IS NULL OR age >= $1) AND
			($2 IS NULL OR age <= $2) AND
			($3 IS NULL OR created_at >= $3) AND
			($4 IS NULL OR created_at <= $4)
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
