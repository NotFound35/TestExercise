package controller

import (
	"awesomeProject/internal/userservice"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	userService *userservice.UserService
	log         *zap.Logger
}

// конструктор handler
func NewHandler(userService *userservice.UserService, log *zap.Logger) *Handler {
	return &Handler{userService: userService, log: log}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) { //interface - чтобы принимать любые данные //code - 200, 404 ...
	w.Header().Set("Content-Type", "application/json")         //что бы ответ был в JSON
	w.WriteHeader(code)                                        //задает HTTP-статус код
	if err := json.NewEncoder(w).Encode(payload); err != nil { //все в JSON
		http.Error(w, "ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, Response{Errors: []string{message}})
}

// w http.ResponseWriter - куда пишется ответ
//Создает структуру Response, где поле Errors содержит переданное сообщение в виде массива строк (даже если ошибка одна).
