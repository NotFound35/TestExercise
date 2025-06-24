package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type GetUserResponse struct {
	Users  []models.User `json:"users"`
	Errors []string      `json:"errors,omitempty"`
}

type GetUserParams struct {
	firstName string
	lastName  string
	age       *int
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.GetUserHandler"
	ctx, cansel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cansel()

	var req Request
	if err := ValidateJSONBody(r, &req, h.Log, op); err != nil {
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := r.URL.Query()
	var params GetUserParams

	params.firstName = strings.TrimSpace(query.Get("first_name"))
	params.lastName = strings.TrimSpace(query.Get("last_name"))
	if ageStr := query.Get("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err != nil {
			h.Log.Error("error strconv",
				zap.String("op", op),
				zap.Error(err),
			)
			responseWithError(w, http.StatusBadRequest, "incorrect age")
			return
		}
		params.age = &age
	}

	//todo написать валидацию JSON

	if err := h.ValidateGetUserParams(params); err != nil {
		h.Log.Error(" ",
			zap.String("op", op),
			zap.Any("params", params),
			zap.Error(err))
		responseWithError(w, http.StatusBadRequest, "validation error")
		return
	}

	users, err := h.UserService.GetUser(ctx, params.firstName, params.lastName, *params.age)
	if err != nil {
		h.Log.Error("error getting user",
			zap.String("op", op),
			zap.Error(err))
		responseWithError(w, http.StatusInternalServerError, "couldnt get user")
		return
	}

	h.Log.Info("user is getting")
	responseWithJson(w, http.StatusOK, map[string]interface{}{
		"answer": users,
	})
}

func (h *Handler) ValidateGetUserParams(params GetUserParams) error {
	var validationErrors []string

	if params.firstName == "" {
		validationErrors = append(validationErrors, "enter your first name")
	}

	if params.lastName == "" {
		validationErrors = append(validationErrors, "enter your last name")
	}

	if params.age != nil {
		if *params.age < 0 {
			validationErrors = append(validationErrors, "age must be greater than zero")
		} else if *params.age > 120 {
			validationErrors = append(validationErrors, "age must be less than 120")
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf(strings.Join(validationErrors, ", "))
	}

	return nil
}
