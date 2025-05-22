package userservice

import (
	"awesomeProject/internal/domain/models"
	"fmt"
)

func (u *UserService) UserSave(user *models.User) error {
	const op = "SaveUser"

	err := u.db.SaveUser(user)
	if err != nil {
		return fmt.Errorf("op: %v, ошибка сохранения юзера %w", op, err)
	}
	return nil
}

func (u *UserService) UserGet(firstName, lastName string, age int) ([]models.User, error) {
	return u.db.GetUserPostgreSQL(firstName, lastName, age)
}
