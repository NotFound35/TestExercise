package repo_test

import (
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/repository/postgresql"
	"context"
	"database/sql"
	"testing"
	_ "time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostgreSQL_SaveUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockSetup   func()
		expectedErr bool
	}{
		{
			name: "success",
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(sqlmock.AnyArg(), "Test", "Test", 50, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: false,
		},
		{
			name: "db error",
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users").
					WillReturnError(sql.ErrConnDone)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			repo := &postgresql.PostgreSQL{
				Db:     db,
				Logger: zap.NewNop(),
			}

			err := repo.SaveUser(context.Background(), &models.User{
				FirstName: "Test",
				LastName:  "Test",
				Age:       50,
			})

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
