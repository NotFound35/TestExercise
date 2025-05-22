package postgresql

import (
	"awesomeProject/internal/domain/models"
	"fmt"
)

func (p *PostgreSQL) GetUserPostgreSQL(firstName, lastName string, age int) ([]models.User, error) {
	const op = "GetUserPostgreSQL"

	//инициализация базового запроса и переменных
	query := "SELECT id, first_name, last_name, age FROM users WHERE 1=1"
	var search []interface{} //для хранения параметров запроса
	var answer []models.User //для хранения результатов
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

	//выполнение SQL-запроса
	result, err := p.db.Query(query, search...)
	if err != nil {
		return nil, fmt.Errorf("op: %s, %w", op, err)
	}
	defer result.Close()

	for result.Next() {
		var user models.User
		if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("op: %s, %w", op, err)
		}
		answer = append(answer, user)
	}

	if err := result.Err(); err != nil {
		return nil, fmt.Errorf("op: %s, %w", op, err)
	}

	return answer, nil
}
