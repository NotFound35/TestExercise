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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{Errors: []string{message}})
}
