package mocks

import (
	"awesomeProject/internal/domain/models"
	"context"
	"github.com/stretchr/testify/mock"
)

type UserDB struct {
	mock.Mock
}

func (m *UserDB) SaveUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserDB) GetUserPostgreSQL(ctx context.Context, firstName, lastName string, age int) ([]models.User, error) {
	args := m.Called(ctx, firstName, lastName, age)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *UserDB) ListUsersPostgreSQL(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error) {
	args := m.Called(ctx, minAge, maxAge, startDate, endDate)
	return args.Get(0).([]models.User), args.Error(1)
}
