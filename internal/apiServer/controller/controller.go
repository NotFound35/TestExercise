package controller

import (
	"awesomeProject/internal/userservice"
	"go.uber.org/zap"
)

type Handler struct {
	userService *userservice.UserService
	log         *zap.Logger
}

// конструктор handler
func NewHandler(userService *userservice.UserService, log *zap.Logger) *Handler {
	return &Handler{userService: userService, log: log}
}
