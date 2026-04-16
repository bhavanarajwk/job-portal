package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	uploadDir    = "uploads/resumes"
	maxFileSize  = 5 << 20 // 5 MB
)

// AllowedResumeTypes lists accepted MIME types for resume uploads
var AllowedResumeTypes = map[string]bool{
	"application/pdf": true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// SaveResume validates and saves an uploaded resume file.
// Returns the relative file path to store in the DB.
func SaveResume(file *multipart.FileHeader) (string, error) {
	// Check file size
	if file.Size > maxFileSize {
		return "", fmt.Errorf("file too large: max size is 5MB")
	}

	// Check MIME type
	contentType := file.Header.Get("Content-Type")
	if !AllowedResumeTypes[contentType] {
		return "", fmt.Errorf("invalid file type: only PDF and Word documents are allowed")
	}

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename: <uuid>_<original_name>
	ext := filepath.Ext(file.Filename)
	safeFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	destPath := filepath.Join(uploadDir, safeFilename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Write to destination
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return "", fmt.Errorf("failed to write file: %w", werr)
			}
		}
		if err != nil {
			break
		}
	}

	// Return forward-slash path for cross-platform consistency
	return strings.ReplaceAll(destPath, "\\", "/"), nil
}
