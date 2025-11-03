package handlers

import (
	"encoding/json"
	"errors"
	apperrors "learning/internal/errors"
	"learning/internal/models"
	"learning/internal/services"
	"learning/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.IUserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes registers user-related routes
func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.GetByID).Methods(http.MethodGet)
}

// Create handles user creation requests
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert User to UserResponse to exclude password
	userResponse := models.ToUserResponse(user)
	utils.WriteSuccess(w, http.StatusCreated, userResponse)
}

// GetByID handles user retrieval by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetUserById(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert User to UserResponse to exclude password
	userResponse := models.ToUserResponse(user)
	utils.WriteSuccess(w, http.StatusOK, userResponse)
}

// handleError processes errors and returns appropriate HTTP responses
func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		utils.WriteError(w, appErr.Code, appErr.Message)
		return
	}

	// Handle validation errors
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		utils.WriteError(w, http.StatusBadRequest, "validation failed: "+validationErr.Error())
		return
	}

	// Default to internal server error
	utils.WriteError(w, http.StatusInternalServerError, "internal server error")
}
