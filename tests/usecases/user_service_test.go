package usecases_test

import (
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/userservice"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// todo 2 мок бд
type MockUserDB struct {
	mock.Mock
}

func (m *MockUserDB) SaveUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserDB) GetUserPostgreSQL(ctx context.Context, firstName, lastName string, age int) ([]models.User, error) {
	args := m.Called(ctx, firstName, lastName, age)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserDB) ListUsersPostgreSQL(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error) {
	args := m.Called(ctx, minAge, maxAge, startDate, endDate)
	return args.Get(0).([]models.User), args.Error(1)
}

func TestUserService_Save(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockUserDB)
		expectedErr error
	}{
		{
			name: "success",
			mockSetup: func(m *MockUserDB) {
				m.On("SaveUser", mock.Anything, &models.User{
					FirstName: "Test",
					LastName:  "Test",
					Age:       50,
				}).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "db error",
			mockSetup: func(m *MockUserDB) {
				m.On("SaveUser", mock.Anything, mock.Anything).
					Return(errors.New("db error"))
			},
			expectedErr: errors.New("метод: SaveUser: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockUserDB)
			tt.mockSetup(mockDB)

			service := &userservice.UserService{
				Db:  mockDB,
				Log: zap.NewNop(),
			}

			err := service.UserSave(context.Background(), &models.User{
				FirstName: "Test",
				LastName:  "Test",
				Age:       50,
			})

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockDB.AssertExpectations(t)
		})
	}
}
