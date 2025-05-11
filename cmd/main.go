package main

import (
	"awesomeProject/config"
	"awesomeProject/internal/controller"
	"awesomeProject/internal/httpServer"
	"awesomeProject/internal/repository/postgresql"
	"awesomeProject/internal/userservice"
	"awesomeProject/logger"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New()

	db, _ := postgresql.NewPostgreSQL(cfg, log)

	defer func(db *postgresql.PostgreSQL) {
		err := db.Close()
		if err != nil {
			log.Error("соединение с БД НЕ закрыто", zap.Error(err))
		}
		log.Info("соединение с БД закрыто")
	}(db)

	userService := userservice.NewUserService(db, log)

	result, err := userService.SaveUser()
	if err != nil {
		log.Error("юзер не сохранен - ОШИБКА")
	}
	log.Info(result)

	// создание http обработчиков
	userHandler := controller.NewUserHandler(userService, log)

	// настройка маршрутов
	router := http.NewServeMux()
	userHandler.RegisterRoutes(router)

	// создание http сервера
	server := httpServer.New(cfg)
	log.Info("создание сервера")

	// запуск http сервера
	server.Start(router)

}
