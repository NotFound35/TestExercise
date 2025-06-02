package controllers

import (
	"awesomeProject/internal/userservice"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	UserService *userservice.UserService
	Log         *zap.Logger
}

func NewHandler(userService *userservice.UserService, log *zap.Logger) *Handler {
	return &Handler{UserService: userService, Log: log}
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	responseWithJson(w, code, Response{Errors: []string{message}})
}
