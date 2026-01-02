package user

import (
	"encoding/json"
	"learning/internal/errors"
	"learning/internal/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  *Service
	validate *validator.Validate
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Create(r.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteError(w, appErr.Code, appErr.Message)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, user)
}
