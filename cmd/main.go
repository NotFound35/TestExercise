package main

import (
	"awesomeProject/internal/apiServer/httpServer"
	"awesomeProject/internal/config"
	"awesomeProject/internal/repository/postgresql"
	"awesomeProject/internal/userservice"
	"awesomeProject/logger"
	"go.uber.org/zap"
)

func main() {
	log := logger.New()

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal("failed to load config",
			zap.Error(err))
	}

	db, err := postgresql.NewPostgreSQL(cfg)
	if err != nil {
		log.Fatal("error connecting to database",
			zap.Error(err))
	}

	userService := userservice.NewUserService(db, log)

	defer db.Close()

	server := httpServer.NewServer(userService, log)

	server.Run(cfg)
}
