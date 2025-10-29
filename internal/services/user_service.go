package services

import (
	"context"
	"fmt"
	"learning/internal/models"
	"learning/internal/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository  *repository.UserRepository) *UserService {
	 return &UserService{
		userRepository: userRepository,
	 }
}

func (s *UserService) hashPassword(password string) (string , error){
	 hash, err := bcrypt.GenerateFromPassword([]byte(password) , bcrypt.DefaultCost);
	 if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	 }
	 return string(hash), nil
}

func (s *UserService) validateCreateUserRequests(req *models.CreateUserRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User,error) {

	if err := s.validateCreateUserRequests(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	hashedPassword , err := s.hashPassword(req.Password); if err != nil {
		return nil, fmt.Errorf("error while hashing password %w" , err );
	}

	user , err := s.userRepository.CreateUser(ctx,req,hashedPassword);
	if err != nil {
		return nil , fmt.Errorf("failed to create user %w", err);
	}

	return user, nil;
 
}


func (s * UserService) GetUserById( ctx context.Context, id int) (*models.User, error) {

	if id < 0 {
		 return nil , fmt.Errorf("invalid id %d", id)
	}

	user, err := s.userRepository.GetUserById(ctx,id) ;
	if err != nil {
		 return nil,  fmt.Errorf("error while getting user by id %w" , err)
	}
	return user, nil;

}