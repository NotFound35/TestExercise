package controllers

import (
	"awesomeProject/internal/domain/models"
	"net/http"
	"strconv"
)

type GetUserResponse struct {
	Users  []models.User `json:"users"`
	Errors []string      `json:"errors,omitempty"`
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// парсим параметры напрямую из URL
	firstName := r.URL.Query().Get("first_name") //Query() - парсит query-запрос в map
	lastName := r.URL.Query().Get("last_name")
	ageStr := r.URL.Query().Get("age")

	// конвертируем возраст из string (ageStr) в int (age)
	var age int
	if ageStr != "" {
		age, _ = strconv.Atoi(ageStr)
	}

	users, err := h.userService.UserGet(firstName, lastName, age)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//формирование ответика
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"answer": users,
	})
}
