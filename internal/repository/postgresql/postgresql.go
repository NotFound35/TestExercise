package postgresql

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type UserDB interface {
	SaveUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	SoftDeleteUser(ctx context.Context, user *models.User) error
	GetUserPostgreSQL(ctx context.Context, firstName, lastName string, age int) ([]models.User, error)
	ListUsersPostgreSQL(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error)
	UserUpdate(ctx context.Context, user *models.User) error
}

type PostgreSQL struct {
	Db *sql.DB
}

//func (p *PostgreSQL) UserDelete(ctx context.Context, user *models.User) error {
//	//TODO implement me
//	panic("implement me")
//}

func NewPostgreSQL(cfg *config.Config) (*PostgreSQL, error) {
	const op = "NewPostgreSQL"

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("op %s: err %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("op %s: err %w", op, err)
	}

	dataB := PostgreSQL{Db: db}
	Migrate(&dataB)

	return &dataB, nil
}

func (p *PostgreSQL) Close() error {
	const op = "Close"
	if err := p.Db.Close(); err != nil {
		return fmt.Errorf("op %s: err %w", op, err)
	}
	return nil
}
