package userservice

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
)

func (u *UserService) SaveUser(ctx context.Context, user *models.User) error {
	const op = "SaveUser"

	if err := u.Db.SaveUser(ctx, user); err != nil {
		return fmt.Errorf("метод: %s: %w", op, err)
	}
	u.Log.Info("User created")

	return nil
}

func (u *UserService) GetUser(ctx context.Context, firstName, lastName string, age int) ([]models.User, error) {
	return u.Db.GetUserPostgreSQL(ctx, firstName, lastName, age)
}

func (u *UserService) ListUsers(
	ctx context.Context,
	minAge, maxAge *int,
	startDate, endDate *int64,
) ([]models.User, error) {
	return u.Db.ListUsersPostgreSQL(ctx, minAge, maxAge, startDate, endDate)
}

func (u *UserService) UserDelete(ctx context.Context, user *models.User) error {
	const op = "DeleteUser"
	if err := u.Db.DeleteUser(ctx, user); err != nil {
		return fmt.Errorf("метод: %s: %w", op, err)
	}
	u.Log.Info("User deleted")
	return nil
}
