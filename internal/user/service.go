package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface defines business operations for users
type ServiceInterface interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
}

// Ensure Service implements ServiceInterface
var _ ServiceInterface = (*Service)(nil)

type Service struct {
	repository RepositoryInterface
	validator  *validator.Validate
}

// NewService creates a new user service
func NewService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
		validator:  validator.New(),
	}
}

func (s *Service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (s *Service) validateCreateUserRequests(req *CreateUserRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("required fields cannot be empty")
	}
	return nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	if err := s.validateCreateUserRequests(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password %w", err)
	}

	user, err := s.repository.CreateUser(ctx, req, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %w", err)
	}

	return user, nil
}

// GetUserById retrieves a user by ID
func (s *Service) GetUserById(ctx context.Context, id int) (*User, error) {
	if id < 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	user, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by id %w", err)
	}
	return user, nil
}
