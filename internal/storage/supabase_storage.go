package storage

import (
	"fmt"
	"io"
	"log"

	storage_go "github.com/supabase-community/storage-go"
)

// StorageService defines the contract for our storage operations.
type StorageService interface {
	UploadImage(bucketName string, filePath string, file io.Reader, contentType string) (string, error)
}

type supabaseStorage struct {
	client *storage_go.Client
	apiURL string
}

// NewSupabaseStorage initializes the Supabase client using your env variables
func NewSupabaseStorage(projectURL string, serviceKey string) StorageService {
	// Initialize the storage client with the project URL and service role key
	client := storage_go.NewClient(projectURL, serviceKey, nil)

	return &supabaseStorage{
		client: client,
		apiURL: projectURL,
	}
}

// UploadImage handles streaming the file to Supabase and returning the public URL
func (s *supabaseStorage) UploadImage(bucketName string, filePath string, file io.Reader, contentType string) (string, error) {
	// 1. Upload the file to the Supabase bucket
	// We pass the io.Reader directly, so it streams from memory without touching the local disk
	res, err := s.client.UploadFile(bucketName, filePath, file)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to supabase: %w", err)
	}

	// Just a quick log to confirm it worked during development
	log.Printf("File uploaded successfully to Supabase. Response: %v", res)

	// 2. Construct and return the Public URL
	// Supabase public URLs always follow this exact structure
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.apiURL, bucketName, filePath)

	return publicURL, nil
}