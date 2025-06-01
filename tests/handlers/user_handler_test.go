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

type MockUserService struct {
	mock.Mock
	*userservice.UserService
}

func (m *MockUserService) UserSave(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestSaveUserHandler(t *testing.T) {

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: `{
				"user": {
					"first_name": "Test",
					"last_name": "Test",
					"age": 50
				}
			}`,
			mockSetup: func(m *MockUserService) {
				m.On("UserSave", mock.Anything, &models.User{
					FirstName: "Test",
					LastName:  "Test",
					Age:       50,
				}).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid json",
			requestBody: `{
				"user": {
					"first_name": "Test",
					"last_name": "Test",
					"age": "fifty"
				}
			}`,
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			h := &controllers.Handler{
				UserService: mockService.UserService,
				Log:         zap.NewNop(),
			}

			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Post("/users", h.SaveUserHandler)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
