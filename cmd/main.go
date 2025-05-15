package main

import (
	"awesomeProject/internal/apiServer/httpServer"
	"awesomeProject/internal/config"
	"awesomeProject/internal/repository/postgresql"
	"awesomeProject/internal/userservice"
	"awesomeProject/logger"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New()

	db, err := postgresql.NewPostgreSQL(cfg, log)
	if err != nil {
		fmt.Println("проблема с sql.DB", err)
	}

	userService := userservice.NewUserService(db, log)

	defer db.Close() // написал так, потому что в функции Сlose уже есть проверка

	server := httpServer.NewServer(userService)

	server.StartAndFinish(cfg) //run || ListenAndServe || ...
}
