package userservice

import (
	"awesomeProject/internal/domain/models"
	"go.uber.org/zap"
)

// интерфейс, который определяет контракт для работы с БД
type UserDB interface {
	SaveUser(user *models.User) error
}

// слой бизнес-логики для работы с юзерами
type UserService struct {
	db  UserDB
	log *zap.Logger
}

func NewUserService(db UserDB, log *zap.Logger) *UserService {
	return &UserService{
		db:  db,
		log: log,
	}
}
