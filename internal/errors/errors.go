package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors
var (
	ErrNotFound            = &AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrBadRequest          = &AppError{Code: http.StatusBadRequest, Message: "bad request"}
	ErrUnauthorized        = &AppError{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden           = &AppError{Code: http.StatusForbidden, Message: "forbidden"}
	ErrInternalServer      = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrValidationFailed    = &AppError{Code: http.StatusBadRequest, Message: "validation failed"}
	ErrConflict            = &AppError{Code: http.StatusConflict, Message: "resource conflict"}
	ErrUnprocessableEntity = &AppError{Code: http.StatusUnprocessableEntity, Message: "unprocessable entity"}
)

// Wrap wraps an error with an AppError
func Wrap(err error, appErr *AppError) *AppError {
	return &AppError{
		Code:    appErr.Code,
		Message: appErr.Message,
		Err:     err,
	}
}

// WrapWithMessage wraps an error with a custom message
func WrapWithMessage(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
