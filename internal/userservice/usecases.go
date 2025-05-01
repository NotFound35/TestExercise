package userservice

import (
	"awesomeProject/internal/domain/models"
	"time"
)

func (u *UserService) SaveUser() (string, error) { // прокинуть сюда модель *models.User, когда буду писать API - передавать его сверху
	user := &models.User{
		ID:            "id_" + time.Now().Format("20060102150405"),
		FirstName:     "Артем",
		LastName:      "Шерр",
		Age:           20,
		RecordingDate: time.Now().Unix(),
	}

	err := u.db.SaveUser(user)
	if err != nil {
		return "не удалось сохранить пользователя", err
	}

	return user.ID, nil
}
