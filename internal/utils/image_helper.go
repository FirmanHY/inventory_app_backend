package utils

import (
	"errors"
	"fmt"
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/pkg/firebase"
	"mime/multipart"
	"path/filepath"
)

func ValidateAndUploadImage(file *multipart.FileHeader, folder string) (string, map[string]string, error) {
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return "", map[string]string{
			"image": constants.MsgImageAllowedFormat,
		}, errors.New("invalid_format")
	}

	if file.Size > 5<<20 { // 5MB
		return "", map[string]string{
			"image": constants.MsgImageAllowedSizes,
		}, errors.New("too_large")
	}

	url, err := firebase.UploadFile(file, folder)
	if err != nil {
		return "", nil, fmt.Errorf("upload_failed: %w", err)
	}

	return url, nil, nil
}
