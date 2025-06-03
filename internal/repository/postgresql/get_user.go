package postgresql

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
)

func (p *PostgreSQL) GetUserPostgreSQL(ctx context.Context, firstName, lastName string, age int) ([]models.User, error) {
	const op = "GetUserPostgreSQL"

	query := "SELECT id, first_name, last_name, age FROM users WHERE 1=1"
	var search []interface{}
	var answer []models.User
	paramCnt := 1

	if firstName != "" {
		query += fmt.Sprintf(" AND first_name ILIKE $%d", paramCnt)
		search = append(search, firstName+"%")
		paramCnt++
	}

	if lastName != "" {
		query += fmt.Sprintf(" AND last_name ILIKE $%d", paramCnt)
		search = append(search, lastName+"%")
		paramCnt++
	}

	if age > 0 {
		query += fmt.Sprintf(" AND age = $%d", paramCnt)
		search = append(search, age)
	}

	rows, err := p.Db.QueryContext(ctx, query, search...)
	if err != nil {
		return nil, fmt.Errorf("op: %s, %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("op: %s, %w", op, err)
		}
		answer = append(answer, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("op: %s, %w", op, err)
	}

	return answer, nil
}
