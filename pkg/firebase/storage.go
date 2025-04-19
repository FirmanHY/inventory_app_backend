package firebase

import (
	"context"
	"fmt"
	"inventory_app_backend/internal/config"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"

	gcsStorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	firebaseStorage "firebase.google.com/go/v4/storage"

	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var StorageClient *firebaseStorage.Client

func InitializeStorage() error {
	credPath := config.Get("FIREBASE_CREDENTIALS_FILES")

	credJSON, err := os.ReadFile(credPath)
	if err != nil {
		return fmt.Errorf("gagal membaca file credentials: %w", err)
	}

	opt := option.WithCredentialsJSON(credJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	StorageClient, err = app.Storage(context.Background())
	return err
}

func UploadFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	ctx := context.Background()
	bucketName := config.Get("FIREBASE_BUCKET_NAME")

	// Generate unique filename
	extension := filepath.Ext(fileHeader.Filename)
	fileName := uuid.New().String() + extension
	contentType := mime.TypeByExtension(extension)

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Upload file
	bucket, err := StorageClient.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	obj := bucket.Object(folder + "/" + fileName)
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err = io.Copy(wc, file); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	// Buat objek bisa diakses publik
	if err := obj.ACL().Set(ctx, gcsStorage.AllUsers, gcsStorage.RoleReader); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", bucketName, folder, fileName)
	return publicURL, nil
}
