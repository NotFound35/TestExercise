package usecases_test

import (
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/userservice"
	"awesomeProject/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestUserService_Save(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*mocks.UserDB)
		expectedErr error
	}{
		{
			name: "success",
			mockSetup: func(m *mocks.UserDB) {
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
			mockSetup: func(m *mocks.UserDB) {
				m.On("SaveUser", mock.Anything, mock.Anything).
					Return(errors.New("db error"))
			},
			expectedErr: errors.New("метод: SaveUser: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.UserDB)
			tt.mockSetup(mockDB)

			service := &userservice.UserService{
				Db:  mockDB,
				Log: zap.NewNop(),
			}

			err := service.SaveUser(context.Background(), &models.User{
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
