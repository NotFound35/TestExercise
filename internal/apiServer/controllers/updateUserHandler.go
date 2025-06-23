package controllers

import (
	"awesomeProject/internal/domain/models"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Handler.UpdateUser"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	userID, err := uuid.Parse(idStr)
	if err != nil {
		h.Log.Error("invalid user id", zap.String("op", op), zap.Error(err))
		responseWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.Log.Error("error decoding body", zap.String("op", op), zap.Error(err))
		responseWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	user.ID = userID

	if err := h.UserService.UserUpdate(ctx, &user); err != nil {
		h.Log.Error("error updating user", zap.String("op", op), zap.Error(err))
		responseWithError(w, http.StatusInternalServerError, "couldnt update user")
		return
	}

	responseWithJson(w, http.StatusOK, user)

}
