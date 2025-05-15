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
	User *models.User `json:"user"`
}

func (h *Handler) SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	var req Request

	//нужно, чтобы JSON был такой же, как и Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неправильный JSON - ", http.StatusBadRequest)
		return //иначе будет использоваться в проге неправильный req
	}
	fmt.Println("правильный JSON")

	if err := Validation(req); err != nil {
		http.Error(w, "не пройдена проверка полей", http.StatusBadRequest)
		return
	}
	fmt.Println("проверка полей пройдена")

	// вызов метода из биз. лог. для сохр. юзера
	_, err := h.userService.SaveUser(req.User) //_ - потому что мне насрать на ID
	if err != nil {
		http.Error(w, "юзер не сохранен - ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated) //ответ, что все ОК
	//fmt.Println("все ОК")
}

func Validation(req Request) error {
	if req.User == nil {
		return errors.New("пустой юзер")
	}

	//Ошибки полей сущности User
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
