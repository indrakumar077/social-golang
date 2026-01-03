package utils

import (
	"io"
	"net/http"
	"strings"
)

// MultipartFile represents an uploaded file
type MultipartFile struct {
	File        io.Reader
	Filename    string
	ContentType string
	FieldName   string
	Size        int64
}

// FileFieldConfig configures a file field
type FileFieldConfig struct {
	FieldName    string   // Form field name (e.g., "photo", "document")
	MaxSize      int64    // Max file size in bytes (0 = no limit)
	AllowedTypes []string // Allowed MIME types (empty = all types allowed)
	Required     bool     // Is this file required?
}

// ParseMultipartRequestOptions configures how to parse multipart requests
type ParseMultipartRequestOptions struct {
	MaxMemory      int64             // Max memory for form parsing (default: 10MB)
	DataFieldName  string            // JSON data field name (default: "data")
	FileFields     []FileFieldConfig // File field configurations
	RequiredFields []string          // Required form fields (non-file)
}

// ParseMultipartRequest parses a request that can be either JSON or multipart/form-data
func ParseMultipartRequest(r *http.Request, opts *ParseMultipartRequestOptions) (jsonData []byte, files map[string]*MultipartFile, err error) {
	// Set defaults
	if opts == nil {
		opts = &ParseMultipartRequestOptions{}
	}
	if opts.MaxMemory == 0 {
		opts.MaxMemory = 10 << 20 // 10MB default
	}
	if opts.DataFieldName == "" {
		opts.DataFieldName = "data"
	}

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form
		if err = r.ParseMultipartForm(opts.MaxMemory); err != nil {
			return nil, nil, &RequestError{
				Message:    "failed to parse form",
				StatusCode: http.StatusBadRequest,
			}
		}

		// Check required non-file fields
		for _, field := range opts.RequiredFields {
			if r.FormValue(field) == "" {
				return nil, nil, &RequestError{
					Message:    field + " field is required",
					StatusCode: http.StatusBadRequest,
				}
			}
		}

		// Get JSON data from form
		dataStr := r.FormValue(opts.DataFieldName)
		if dataStr == "" {
			return nil, nil, &RequestError{
				Message:    opts.DataFieldName + " field is required",
				StatusCode: http.StatusBadRequest,
			}
		}
		jsonData = []byte(dataStr)

		// Extract and validate files
		files = make(map[string]*MultipartFile)
		for _, fileConfig := range opts.FileFields {
			fileHandle, header, fileErr := r.FormFile(fileConfig.FieldName)

			if fileErr != nil {
				// File not provided
				if fileConfig.Required {
					return nil, nil, &RequestError{
						Message:    fileConfig.FieldName + " file is required",
						StatusCode: http.StatusBadRequest,
					}
				}
				continue // Optional file not provided, skip
			}

			// File provided - validate it
			fileSize := header.Size

			// Validate file size
			if fileConfig.MaxSize > 0 && fileSize > fileConfig.MaxSize {
				return nil, nil, &RequestError{
					Message:    fileConfig.FieldName + " file size exceeds maximum allowed size",
					StatusCode: http.StatusBadRequest,
				}
			}

			// Validate file type
			contentType := header.Header.Get("Content-Type")
			if len(fileConfig.AllowedTypes) > 0 {
				if !isAllowedType(contentType, fileConfig.AllowedTypes) {
					return nil, nil, &RequestError{
						Message:    fileConfig.FieldName + " file type not allowed. Allowed types: " + strings.Join(fileConfig.AllowedTypes, ", "),
						StatusCode: http.StatusBadRequest,
					}
				}
			}

			// File is valid, add to map
			files[fileConfig.FieldName] = &MultipartFile{
				File:        fileHandle,
				Filename:    header.Filename,
				ContentType: contentType,
				FieldName:   fileConfig.FieldName,
				Size:        fileSize,
			}
		}

		return jsonData, files, nil
	}

	// Handle JSON request

	// if err != nil {
	// 	return nil, nil, &RequestError{
	// 		Message:    "failed to read request body",
	// 		StatusCode: http.StatusBadRequest,
	// 	}
	// }

	return jsonData, nil, nil
}

// isAllowedType checks if content type is in allowed list
func isAllowedType(contentType string, allowedTypes []string) bool {
	contentType = strings.ToLower(contentType)
	for _, allowed := range allowedTypes {
		if contentType == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}

// RequestError represents a request parsing error
type RequestError struct {
	Message    string
	StatusCode int
}

func (e *RequestError) Error() string {
	return e.Message
}
