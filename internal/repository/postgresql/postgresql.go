package postgresql

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/domain/models"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type UserDB interface {
	SaveUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	GetUserPostgreSQL(ctx context.Context, firstName, lastName string, age int) ([]models.User, error)
	ListUsersPostgreSQL(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error)
}

type PostgreSQL struct {
	Db     *sql.DB
	Logger *zap.Logger
}

//func (p *PostgreSQL) UserDelete(ctx context.Context, user *models.User) error {
//	//TODO implement me
//	panic("implement me")
//}

func NewPostgreSQL(cfg *config.Config, logger *zap.Logger) (*PostgreSQL, error) {
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
		return nil, fmt.Errorf("функция %v: %v", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("функция %v: %v ", op, err)
	}

	logger.Info("успешный коннект с PostgreSQL")

	dataB := PostgreSQL{Db: db, Logger: logger}
	Migrate(&dataB)

	return &dataB, nil
}

func (p *PostgreSQL) Close() {
	const op = "Close"
	if err := p.Db.Close(); err != nil {
		//todo log
		p.Logger.Error("DB close error!!!",
			zap.Error(fmt.Errorf("метод %v: %v", op, err)))
	}
	p.Logger.Info("соединения с БД закрыто - УСПЕХ!!!")
}
