package userservice

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
)

func (u *UserService) SaveUser(ctx context.Context, user *models.User) error {
	const op = "SaveUser"

	if err := u.Db.SaveUser(ctx, user); err != nil {
		return fmt.Errorf("op %s: err %w", op, err)
	}
	u.Log.Info("User created")

	return nil
}

func (u *UserService) GetUser(ctx context.Context, firstName, lastName string, age int) ([]models.User, error) {
	const op = "GetUser"
	users, err := u.Db.GetUserPostgreSQL(ctx, firstName, lastName, age)
	if err != nil {
		return nil, fmt.Errorf("op %s: err %w", op, err)
	}
	return users, nil
}

func (u *UserService) ListUsers(ctx context.Context, minAge, maxAge *int, startDate, endDate *int64) ([]models.User, error) {
	const op = "ListUser"
	users, err := u.Db.ListUsersPostgreSQL(ctx, minAge, maxAge, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("op %s: err %w", op, err)
	}
	return users, nil
}

func (u *UserService) UserDelete(ctx context.Context, user *models.User) error {
	const op = "DeleteUser"
	if err := u.Db.DeleteUser(ctx, user); err != nil {
		return fmt.Errorf("op %s: err %w", op, err)
	}
	u.Log.Info("User deleted")
	return nil
}

func (u *UserService) SoftUserDelete(ctx context.Context, user *models.User) error {
	const op = "SoftUserDelete"
	if err := u.Db.SoftDeleteUser(ctx, user); err != nil {
		return fmt.Errorf("op %s: err %w", op, err)
	}
	u.Log.Info("User soft-deleted")
	return nil
}

func (u *UserService) UserUpdate(ctx context.Context, user *models.User) error {
	const op = "UserService.UserUpdate"
	if err := u.Db.UserUpdate(ctx, user); err != nil {
		return fmt.Errorf("op %s: err %w", op, err)
	}
	u.Log.Info("User updated")
	return nil
}
