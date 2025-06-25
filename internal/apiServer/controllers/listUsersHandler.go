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

type ListUsersResponse struct {
	Users  []models.User `json:"users"`
	Errors []string      `json:"errors,omitempty"`
}

type ListUsersParams struct {
	MinAge    *int
	MaxAge    *int
	StartDate *int64
	EndDate   *int64
}

func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.ListUsersHandler"
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

	var params ListUsersParams
	query := r.URL.Query()

	if minAgeStr := query.Get("min_age"); minAgeStr != "" {
		if age, err := strconv.Atoi(minAgeStr); err == nil {
			params.MinAge = &age
		} else {
			h.Log.Error("invalid min_age",
				zap.String("op", op),
				zap.String("value", minAgeStr),
				zap.Error(err))
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if maxAgeStr := query.Get("max_age"); maxAgeStr != "" {
		if age, err := strconv.Atoi(maxAgeStr); err == nil {
			params.MaxAge = &age
		} else {
			h.Log.Error(" ",
				zap.String("op", op),
				zap.String("value", maxAgeStr),
				zap.Error(err))
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if startDateStr := query.Get("start_date"); startDateStr != "" {
		if date, err := strconv.ParseInt(startDateStr, 10, 64); err == nil {
			params.StartDate = &date
		} else {
			h.Log.Error("invalid start_date",
				zap.String("op", op),
				zap.String("value", startDateStr),
				zap.Error(err))
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if endDateStr := query.Get("end_date"); endDateStr != "" {
		if date, err := strconv.ParseInt(endDateStr, 10, 64); err == nil {
			params.EndDate = &date
		} else {
			h.Log.Error("invalid end_date",
				zap.String("op", op),
				zap.String("value", endDateStr),
				zap.Error(err))
			responseWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	if err := h.ValidationListUsers(params); err != nil {
		h.Log.Error("validation error",
			zap.String("op", op),
			zap.Any("params", params),
			zap.Error(err))
		responseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.UserService.ListUsers(ctx, params.MinAge, params.MaxAge, params.StartDate, params.EndDate)
	if err != nil {
		h.Log.Error("error listing users",
			zap.String("op", op),
			zap.Error(err))
		responseWithError(w, http.StatusInternalServerError, "couldnt list users")
		return
	}

	responseWithJson(w, http.StatusOK, ListUsersResponse{
		Users: users,
	})
}

func (h *Handler) ValidationListUsers(params ListUsersParams) error {
	var validationErrors []string

	if params.MinAge != nil && *params.MinAge < 0 {
		validationErrors = append(validationErrors, "min age must not be negative")
	}

	if params.MaxAge != nil && *params.MaxAge < 0 {
		validationErrors = append(validationErrors, "max age must not be negative")
	}

	if params.MinAge != nil && params.MaxAge != nil && *params.MinAge > *params.MaxAge {
		validationErrors = append(validationErrors, "min age must not be greater than max age")
	}

	if params.StartDate != nil && *params.StartDate < 0 {
		validationErrors = append(validationErrors, "start date must not be negative")
	}

	if params.EndDate != nil && *params.EndDate < 0 {
		validationErrors = append(validationErrors, "end date must not be negative")
	}

	if params.StartDate != nil && params.EndDate != nil && *params.StartDate > *params.EndDate {
		validationErrors = append(validationErrors, "start date must not be greater than end date")
	}

	if len(validationErrors) > 0 { //если есть ошибки
		return fmt.Errorf(strings.Join(validationErrors, ", "))
	}

	return nil
}
