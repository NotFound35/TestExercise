package userservice

import (
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/repository/postgresql"
	"context"

	"go.uber.org/zap"
)

type IUserService interface {
	GetUser(ctx context.Context, firstName, lastName string, age int) ([]models.User, error)
	SaveUser(ctx context.Context, user *models.User) error
	ListUsers(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error)
}

type UserService struct {
	Db  postgresql.UserDB
	Log *zap.Logger
}

func NewUserService(db postgresql.UserDB, log *zap.Logger) *UserService {
	return &UserService{
		Db:  db,
		Log: log,
	}
}
