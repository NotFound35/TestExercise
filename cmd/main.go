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
		fmt.Println("проблема с соединения с бд", err)
	}

	userService := userservice.NewUserService(db, log)

	defer db.Close()

	server := httpServer.NewServer(userService, log)

	server.Run(cfg)
}
