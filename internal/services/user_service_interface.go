package services

import (
	"context"

	"learning/internal/models"
)

// IUserService defines the interface for user service operations
type IUserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
}
