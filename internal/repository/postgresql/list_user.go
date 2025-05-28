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
	` //NULL - если не указано то пофиг

	//запрос в БД с параметрами
	result, err := p.db.QueryContext(ctx, query, minAge, maxAge, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("op: %s, %w", op, err)
	}
	defer result.Close()

	//чтение результата
	var users []models.User
	for result.Next() { //перебор строк
		var user models.User
		err := result.Scan( //копирование строк из БД в структуру User
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.RecordingDate,
		)
		if err != nil {
			return nil, fmt.Errorf("op: %s, %w", op, err)
		}
		users = append(users, user) //добавление пользователей в итоговой список
	}

	return users, nil
}
