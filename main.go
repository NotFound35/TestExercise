package main

import (
	"awesomeProject/config"
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/infrastructure/postgresql"
	"awesomeProject/logger"
	"fmt"
	"time"
)

func main() {
	// 1. Загрузка конфига
	cfg := config.MustLoad()

	// 2. Инициализация логгера
	log := logger.New()
	defer log.Sync()

	// 3. Инициализация подключения к БД с логгером
	db, _ := postgresql.NewPostgreSQL(cfg, log)
	fmt.Println("коннект с БД произошел!!!")

	defer db.Close()
	fmt.Println("")

	// 4. Создание таблиц
	if err := db.CreateTables(); err != nil {
		fmt.Println("НЕ создалась таблица")
	}

	// 5. Создание тестового пользователя
	user := &models.User{
		ID:            "id_" + time.Now().Format("20060102150405"),
		FirstName:     "Артем",
		LastName:      "Арефьев",
		Age:           27,
		RecordingDate: time.Now().Unix(),
	}

	// 6. Сохранение пользователя
	if err := db.SaveUser(user); err != nil {
		fmt.Println("юзер НЕ сохранен")
	}

	fmt.Println("юзер сохранен!!!")
}

// использовать только один объект с 19ой строчки (соединения с) БД, второй удалить
