package controller

import (
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/userservice"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type UserHandler struct { //главный обработчик http запросов
	userService *userservice.UserService
	logger      *zap.Logger
}

// создает экземпляр обработчика - удобно для тестирования
func NewUserHandler(userService *userservice.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// DTO - Data Transfer Objects
// UserRequest структура входящего запроса
type UserRequest struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Age           int    `json:"age"`
	RecordingDate int64  `json:"recording_date,omitempty"`
}

// UserResponse структура ответа
type UserResponse struct {
	ID            string `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Age           int    `json:"age"`
	RecordingDate int64  `json:"recording_date"`
}

func (h *UserHandler) SaveUser(w http.ResponseWriter, r *http.Request) {
	//проверка метода (должен быть POST)
	if r.Method != http.MethodPost {
		http.Error(w, "метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	//Читает тело запроса и преобразует JSON в структуру UserRequest
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { //r.Body - поток с данными запроса; json.NewDecoder - декодирует JSON построчно
		h.logger.Error("запрос не расшифрован", zap.Error(err))
		http.Error(w, "неправильный запрос", http.StatusBadRequest)
		return
	}

	//валидация данных сюда приделали
	if err := validateUserRequest(req); err != nil {
		h.logger.Error("проверка не пройдена", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//преобразование в модель
	//Преобразует DTO (UserRequest) в доменную модель (models.User)
	user := &models.User{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Age:           req.Age,
		RecordingDate: req.RecordingDate,
	}

	//формирование ответа - создание объекта для ответа клиенту
	response := UserResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Age:           user.Age,
		RecordingDate: user.RecordingDate,
	}

	//отправка ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("ответ не закодирован", zap.Error(err))
	}
	//1. устанавливает заголовок Content-Type: application/json
	//2. статус 201 Created (успешное создание ресурса)
	//3. кодирует ответ в JSON и отправляет
}

// регистрация HTTP маршруты для обработки запросов, связанных с пользователями
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", h.userRoutes)
}

// метод-маршрутизатор, который определяет, какой обработчик вызвать в зависимости от типа http запроса
func (h *UserHandler) userRoutes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.SaveUser(w, r)
	default:
		http.Error(w, "неправильный метод", http.StatusMethodNotAllowed)
	}
}

// валидация - проверка всех полей модели
func validateUserRequest(req UserRequest) error {
	if strings.TrimSpace(req.FirstName) == "" {
		return errors.New("укажите имя")
	}
	if len(req.FirstName) < 2 {
		return errors.New("нужно больше 2х символов")
	}
	if strings.TrimSpace(req.LastName) == "" {
		return errors.New("укажите фамилию")
	}
	if req.Age <= 0 || req.Age > 100 {
		return errors.New("возраст от 0 до 100")
	}
	return nil
}
