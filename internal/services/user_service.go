package services

import (
	"context"
	"fmt"
	"learning/internal/models"
	"learning/internal/repository"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// Ensure UserService implements IUserService interface
var _ IUserService = (*UserService)(nil)

type UserService struct {
	userRepository repository.IUserRepository
	validator      *validator.Validate
}

// NewUserService creates a new user service
func NewUserService(userRepository repository.IUserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
		validator:      validator.New(),
	}
}

func (s *UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (s *UserService) validateCreateUserRequests(req *models.CreateUserRequest) error {
	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("required fields cannot be empty")
	}
	return nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	if err := s.validateCreateUserRequests(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password %w", err)
	}

	user, err := s.userRepository.CreateUser(ctx, req, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %w", err)
	}

	return user, nil
}

// GetUserById retrieves a user by ID
func (s *UserService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	if id < 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by id %w", err)
	}
	return user, nil
}
