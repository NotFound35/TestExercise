package userservice

import (
	"awesomeProject/internal/domain/models"
	"errors"
	"fmt"
)

func (u *UserService) SaveUser(user *models.User) (string, error) { // прокинуть сюда модель *models.User, когда буду писать API - передавать его сверху
	if u.db == nil {
		return "", errors.New("БД не init")
	}

	fmt.Println(user)

	err := u.db.SaveUser(user)
	if err != nil {
		return "не удалось сохранить пользователя", err
	}
	return user.ID, nil
}
