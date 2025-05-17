package controller

import (
	"awesomeProject/internal/domain/models"
	"awesomeProject/internal/userservice"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Handler struct {
	userService *userservice.UserService
}

// конструктор handler
func NewHandler(userService *userservice.UserService) *Handler {
	return &Handler{userService: userService}
}

type Request struct {
	User struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
	} `json:"user"`
}

type Response struct {
	Message string `json:"message"`
}

func (h *Handler) SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.SaveUserHandler"
	var req Request

	//нужно, чтобы JSON был такой же, как и Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//http.Error(w, "неправильный JSON - ", http.StatusBadRequest)
		fmt.Errorf("метод %v: %v", op, err)
		return //иначе будет использоваться в проге неправильный req
	}
	fmt.Println("правильный JSON")

	if err := Validation(req); err != nil {
		//http.Error(w, "не пройдена проверка полей", http.StatusBadRequest)
		fmt.Errorf("метод %v: %v", op, err)
		return
	}
	fmt.Println("проверка полей пройдена")

	// Создаем пользователя с генерацией ID
	user := models.NewUser(
		req.User.FirstName,
		req.User.LastName,
		req.User.Age,
	)

	//метод из бизнес логики
	_, err := h.userService.SaveUser(user)
	if err != nil {
		//http.Error(w, "юзер не сохранен - ", http.StatusInternalServerError)
		fmt.Errorf("метод %v: %v", op, err)
		return
	}

	response := Response{
		Message: fmt.Sprintf("пользователь %s успешно сохранен", req.User.FirstName),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //ответ, что все ОК
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Errorf("метод %v: %v - ошибка при формировании ответа", op, err)
		fmt.Errorf("метод %v: %v - юзер не сохранен", op, err)
	}
}

func Validation(req Request) error {
	if req.User.FirstName == "" {
		return errors.New("нет имени")
	}
	if req.User.LastName == "" {
		return errors.New("нет фамилии")
	}
	if req.User.Age == 0 {
		return errors.New("нет возраста")
	}

	return nil
}
