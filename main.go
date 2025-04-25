package main

import (
	"awesomeProject/config"
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/repository/postgresql"
	"awesomeProject/logger"
	"go.uber.org/zap"
	"time"
)

func main() {
	// 1. Загрузка конфига
	cfg := config.MustLoad()

	// 2. Инициализация логгера
	log := logger.New()
	defer log.Sync()

	// 3. Инициализация подключения к БД с логгером
	db, err := postgresql.NewPostgreSQL(cfg, log)
	if err != nil {
		log.Error("не произошло коннекта с БД", zap.Error(err))
	}
	log.Info("коннект с БД произошел!!!")

	defer func() {
		db.Close()
		log.Info("соединение с БД закрыто!!!")
	}()

	// 4. Создание таблиц
	postgresql.Migrate(db)()
	log.Info("миграции применены!!!")

	// 5. Создание тестового пользователя
	user := &models.User{
		ID:            "id_" + time.Now().Format("20060102150405"),
		FirstName:     "Артем",
		LastName:      "Арефьев",
		Age:           27,
		RecordingDate: time.Now().Unix(),
	}

	// 6. Сохранение пользователя
	if err = db.SaveUser(user); err != nil {
		log.Error("юзер НЕ сохранен", zap.Error(err))
	}

}

// использовать только один объект с 19ой строчки (соединения с) БД, второй удалить
