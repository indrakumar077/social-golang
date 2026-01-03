package user

import (
	"encoding/json"
	"io"
	"learning/internal/errors"
	"learning/internal/utils"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	var photoFile io.Reader
	var photoFilename, photoContentType string

	// Configure multipart parsing - sab kuch bahar se
	opts := &utils.ParseMultipartRequestOptions{
		MaxMemory:     10 << 20, // 10MB
		DataFieldName: "data",
		FileFields: []utils.FileFieldConfig{
			{
				FieldName:    "photo",                                          // Field name
				MaxSize:      5 << 20,                                          // 5MB limit
				AllowedTypes: []string{"image/jpeg", "image/png", "image/jpg"}, // Allowed types
				Required:     false,                                            // Optional
			},
		},
	}

	// Parse request
	jsonData, files, err := utils.ParseMultipartRequest(r, opts)
	if err != nil {
		if reqErr, ok := err.(*utils.RequestError); ok {
			utils.WriteError(w, reqErr.StatusCode, reqErr.Message)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(jsonData, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get file (already validated by utility)
	if photoFileData, exists := files["photo"]; exists {
		defer func() {
			if closer, ok := photoFileData.File.(io.Closer); ok {
				closer.Close()
			}
		}()

		photoFile = photoFileData.File
		photoFilename = photoFileData.Filename
		photoContentType = photoFileData.ContentType
	}

	if err := h.validate.Struct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Create(r.Context(), &req, photoFile, photoFilename, photoContentType)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteError(w, appErr.Code, err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, user)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "User id is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadGateway, "Id is not uuid")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {

		if appErr, ok := err.(*errors.AppError); ok {
			utils.WriteError(w, appErr.Code, appErr.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
	}

	utils.WriteSuccess(w, http.StatusOK, user)
}
