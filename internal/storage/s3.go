package storage

import (
	"context"
	"fmt"
	"io"
	"learning/internal/config"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type S3Client struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Client(cfg *config.S3Config) (*S3Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		client: client,
		bucket: cfg.Bucket,
		region: cfg.Region,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// ext := filepath.Ext(filename)
	ext := filepath.Ext(filename)

	uniqueFilename := fmt.Sprintf("users/%s/%s%s", time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	uploader := manager.NewUploader(s.client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(uniqueFilename),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.bucket,
		s.region,
		uniqueFilename,
	)
	return url, nil
}

// DeleteFileByURL deletes a file from S3 using its URL
func (s *S3Client) DeleteFileByURL(ctx context.Context, url string) error {
	// Extract key from URL
	// URL format: https://bucket.s3.region.amazonaws.com/key
	key := s.extractKeyFromURL(url)
	if key == "" {
		return fmt.Errorf("invalid S3 URL format")
	}

	return s.DeleteFile(ctx, key)
}

// DeleteFile deletes a file from S3 using its key
func (s *S3Client) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

// extractKeyFromURL extracts the S3 key from a full URL
func (s *S3Client) extractKeyFromURL(url string) string {
	// URL format: https://bucket.s3.region.amazonaws.com/key
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s.bucket, s.region)
	if !strings.HasPrefix(url, prefix) {
		return ""
	}
	return strings.TrimPrefix(url, prefix)
}
