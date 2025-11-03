package repository

import (
	"context"

	"learning/internal/models"
)

// IUserRepository defines the interface for user repository operations
type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.CreateUserRequest, hashedPassword string) (*models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
}
