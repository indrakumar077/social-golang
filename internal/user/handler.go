package user

import (
	"encoding/json"
	"errors"
	apperrors "learning/internal/errors"
	"learning/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Handler handles user-related HTTP requests
type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new user handler
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers user-related routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.GetByID).Methods(http.MethodGet)
}

// Create handles user creation requests
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert User to UserResponse to exclude password
	userResponse := ToUserResponse(user)
	utils.WriteSuccess(w, http.StatusCreated, userResponse)
}

// GetByID handles user retrieval by ID
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.service.GetUserById(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Convert User to UserResponse to exclude password
	userResponse := ToUserResponse(user)
	utils.WriteSuccess(w, http.StatusOK, userResponse)
}

// handleError processes errors and returns appropriate HTTP responses
func (h *Handler) handleError(w http.ResponseWriter, err error) {
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
