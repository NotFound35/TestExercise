package userservice

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"go.uber.org/zap"
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

func (u *UserService) SoftUserDelete(ctx context.Context, user *models.User) error {
	const op = "SoftUserDelete"
	if err := u.Db.SoftDeleteUser(ctx, user); err != nil {
		return fmt.Errorf("<UNK>: %s: %w", op, err)
	}
	u.Log.Info("User soft-deleted")
	return nil
}

func (u *UserService) UserUpdate(ctx context.Context, user *models.User) error {
	const op = "UserService.UserUpdate"
	if err := u.Db.UserUpdate(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	u.Log.Info("User updated")
	return nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.Db.GetUserByEmail(ctx, email)
	if err != nil {
		u.Log.Error("GetUserByEmail failed", zap.Error(err))
		return nil, err
	}
	return user, nil
}
