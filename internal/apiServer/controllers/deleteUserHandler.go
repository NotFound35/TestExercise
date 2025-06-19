package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// type DeleteRequest struct {
// 	User struct {
// 		FirstName string `json:"first_name"`
// 		LastName  string `json:"last_name"`
// 		Age       int    `json:"age"`
// 	} `json:"user"`
// }

type DeleteResponse struct {
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.DeleteUserHandler"
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		h.Log.Error("пустой id", zap.String("op", op))
		responseWithError(w, http.StatusBadRequest, "пустой id")
		return
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		h.Log.Error("ошибка преобразования (парсинга) id",
			zap.String("op", op),
			zap.String("value", idStr),
			zap.Error(err))
		responseWithError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	user := &models.User{
		ID: userID,
	}

	if err := h.UserService.UserDelete(ctx, user); err != nil {
		h.Log.Error("ошибка удаления пользователя",
			zap.String("op", op),
			zap.Error(err))
		responseWithError(w, http.StatusInternalServerError, "ошибка удаления пользователя")
		return
	}

	h.Log.Info("User deleted")
	responseWithJson(w, http.StatusOK, DeleteResponse{
		Message: "User deleted",
	})

}
