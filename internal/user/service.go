package user

import (
	"context"
	"learning/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.WrapWithMessage(err, 500, "Failed to hash the password")
	}

	user := &User{
		Username:   req.Username,
		Email:      req.Email,
		Phone:      req.Phone,
		Name:       req.Name,
		Password:   string(hashedPassword),
		MiddleName: req.MiddleName,
		Surname:    req.Surname,
		Bio:        req.Bio,
		Active:     true,
	}

	createUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, errors.WrapWithMessage(err, 500, "failed to create user")
	}

	return ToUserReponse(createUser), nil
}
