package handlers_test

import (
	"awesomeProject/internal/apiServer/controllers"
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/userservice"
	_ "awesomeProject/internal/userservice"
	"bytes"
	"context"
	_ "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	_ "time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// todo 1 мок бд
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

func TestSaveUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*MockUserDB)
		expectedStatus int
	}{
		{
			name:        "success",
			requestBody: `{"user": {"first_name":"Test", "last_name":"Test", "age":50}}`,
			mockSetup: func(m *MockUserDB) {
				m.On("SaveUser", mock.Anything, &models.User{
					FirstName: "Test",
					LastName:  "Test",
					Age:       50,
				}).Return(nil)
			},
			expectedStatus: http.StatusCreated,
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

			h := &controllers.Handler{
				UserService: service,
				Log:         zap.NewNop(),
			}

			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/users", h.SaveUserHandler)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDB.AssertExpectations(t)
		})
	}
}
