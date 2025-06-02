package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	var params ListUsersParams

	query := r.URL.Query()

	if minAgeStr := query.Get("min_age"); minAgeStr != "" {
		if age, err := strconv.Atoi(minAgeStr); err == nil {
			params.MinAge = &age
		} else {
			h.Log.Error("ошибка преобразования (парсинга) min_age",
				zap.String("op", op),
				zap.String("value", minAgeStr),
				zap.Error(err))
			respondWithError(w, http.StatusBadRequest, "некорректный min_age")
			return
		}
	}

	if maxAgeStr := query.Get("max_age"); maxAgeStr != "" {
		if age, err := strconv.Atoi(maxAgeStr); err == nil {
			params.MaxAge = &age
		} else {
			h.Log.Error("ошибка преобразования (парсинга) max_age",
				zap.String("op", op),
				zap.String("value", maxAgeStr),
				zap.Error(err))
			respondWithError(w, http.StatusBadRequest, "некорректный max_age")
			return
		}
	}

	if startDateStr := query.Get("start_date"); startDateStr != "" {
		if date, err := strconv.ParseInt(startDateStr, 10, 64); err == nil {
			params.StartDate = &date
		} else {
			h.Log.Error("ошибка преобразования (парсинга) start_date",
				zap.String("op", op),
				zap.String("value", startDateStr),
				zap.Error(err))
			respondWithError(w, http.StatusBadRequest, "некорректный start_date")
			return
		}
	}

	if endDateStr := query.Get("end_date"); endDateStr != "" {
		if date, err := strconv.ParseInt(endDateStr, 10, 64); err == nil {
			params.EndDate = &date
		} else {
			h.Log.Error("ошибка преобразования (парсинга) end_date",
				zap.String("op", op),
				zap.String("value", endDateStr),
				zap.Error(err))
			respondWithError(w, http.StatusBadRequest, "некорректный end_date")
			return
		}
	}

	if err := ValidationListUsers(params); err != nil {
		h.Log.Error("валидация не пройдена", //общий текст ошибки
			zap.String("op", op), //имя текущей операции
			zap.Error(err))       //ошибка, возвращенная валидатором
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.UserService.UsersList(ctx, params.MinAge, params.MaxAge, params.StartDate, params.EndDate)
	if err != nil {
		h.Log.Error("ошибка при получении пользователей",
			zap.String("op", op),
			zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "ошибка сервера")
		return
	}

	respondWithJSON(w, http.StatusOK, ListUsersResponse{
		Users: users,
	})
}

func ValidationListUsers(params ListUsersParams) error {
	var validationErrors []string

	if params.MinAge != nil && *params.MinAge < 0 {
		validationErrors = append(validationErrors, "минимальный возраст не может быть отрицательным")
	}

	if params.MaxAge != nil && *params.MaxAge < 0 {
		validationErrors = append(validationErrors, "максимальный возраст не может быть отрицательным")
	}

	if params.MinAge != nil && params.MaxAge != nil && *params.MinAge > *params.MaxAge {
		validationErrors = append(validationErrors, "минимальный возраст не может быть больше максимального")
	}

	if params.StartDate != nil && *params.StartDate < 0 {
		validationErrors = append(validationErrors, "начальная дата не может быть отрицательной")
	}

	if params.EndDate != nil && *params.EndDate < 0 {
		validationErrors = append(validationErrors, "конечная дата не может быть отрицательной")
	}

	if params.StartDate != nil && params.EndDate != nil && *params.StartDate > *params.EndDate {
		validationErrors = append(validationErrors, "начальная дата не может быть больше конечной")
	}

	if len(validationErrors) > 0 { //если есть ошибки
		return fmt.Errorf(strings.Join(validationErrors, "\n"))
	}

	return nil
}
