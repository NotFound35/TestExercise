package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Request struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

type Response struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func (h *Handler) SaveUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.SaveUserHandler"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var req Request
	if err := validateJSONBody(r, &req, h.Log, op); err != nil {
		h.Log.Error("error validation JSON",
			zap.String("op", op),
			zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Validation(req); err != nil {
		h.Log.Error("validation error",
			zap.String("op", op),
			zap.Any("request", req),
			zap.Error(err),
		)
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
	}

	if err := h.UserService.SaveUser(ctx, user); err != nil {
		h.Log.Error("error saving user",
			zap.String("op", op),
			zap.Error(err),
		)
		responseWithError(w, http.StatusInternalServerError, "couldnt save user")
		return
	}

	h.Log.Info("save user success")
	responseWithJson(w, http.StatusCreated, Response{
		Message: fmt.Sprintf("user %s saved", user.FirstName),
	})
}

func (h *Handler) Validation(req Request) error {
	var validationErrors []string

	firstName := strings.TrimSpace(req.FirstName)
	if firstName == "" {
		validationErrors = append(validationErrors, "enter your first_name")
	}

	lastName := strings.TrimSpace(req.LastName)
	if lastName == "" {
		validationErrors = append(validationErrors, "enter your last_name")
	}

	if req.Age <= 0 {
		validationErrors = append(validationErrors, "age must be greater than zero")
	} else if req.Age > 120 {
		validationErrors = append(validationErrors, "age must be less than 120")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf(strings.Join(validationErrors, ", "))
	}

	return nil
}
