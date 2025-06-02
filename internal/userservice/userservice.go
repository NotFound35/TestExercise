package userservice

import (
	"awesomeProject/internal/domain/models"
	"context"
	"go.uber.org/zap"
)

type UserDB interface {
	SaveUser(ctx context.Context, user *models.User) error
	GetUserPostgreSQL(ctx context.Context, firstName, lastName string, age int) ([]models.User, error)
	ListUsersPostgreSQL(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error)
}

type UserService struct {
	Db  UserDB
	Log *zap.Logger
}

func NewUserService(db UserDB, log *zap.Logger) *UserService {
	return &UserService{
		Db:  db,
		Log: log,
	}
}
