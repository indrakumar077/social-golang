package user

import (
	"context"
	"io"
	"learning/internal/errors"
	"learning/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo     *Repository
	s3Client *storage.S3Client
}

func NewService(repo *Repository, s3Client *storage.S3Client) *Service {
	return &Service{
		repo:     repo,
		s3Client: s3Client,
	}
}

func (s *Service) Create(ctx context.Context, req *CreateUserRequest, photoFile io.Reader, photoFilename, photoContentType string) (*UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.WrapWithMessage(err, 500, "Failed to hash the password")
	}

	var profilePhotoURL *string
	photoUploaded := false

	// Upload photo first (if provided)
	if photoFile != nil && photoFilename != "" {
		url, err := s.s3Client.UploadFile(ctx, photoFile, photoFilename, photoContentType)
		if err != nil {
			return nil, errors.WrapWithMessage(err, 500, "failed to upload profile photo")
		}
		profilePhotoURL = &url
		photoUploaded = true
	}

	user := &User{
		Username:        req.Username,
		Email:           req.Email,
		Phone:           req.Phone,
		Name:            req.Name,
		Password:        string(hashedPassword),
		MiddleName:      req.MiddleName,
		Surname:         req.Surname,
		Bio:             req.Bio,
		ProfilePhotoURL: profilePhotoURL,
		Active:          true,
	}

	// Try to create user
	createUser, err := s.repo.Create(ctx, user)
	if err != nil {
		// If user creation failed and photo was uploaded, delete the photo
		if photoUploaded && profilePhotoURL != nil {
			deleteErr := s.s3Client.DeleteFileByURL(ctx, *profilePhotoURL)
			if deleteErr != nil {
				// Log the error but don't fail the request because of cleanup failure
				// The main error (user creation failure) should be returned
				// In production, you might want to log this to a monitoring system
			}
		}
		return nil, errors.WrapWithMessage(err, 500, "failed to create user")
	}

	return ToUserReponse(createUser), nil
}
