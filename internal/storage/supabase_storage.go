package storage

import (
	"fmt"
	"io"
	"log"
	"strings" // ✅ ADDED: Required to clean the URL

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
	
	// 1. AGGRESSIVE CLEANUP: Remove invisible Windows \r, \n, spaces, quotes, and slashes
	cleanURL := strings.TrimSpace(projectURL)
	cleanURL = strings.ReplaceAll(cleanURL, "\r", "")
	cleanURL = strings.ReplaceAll(cleanURL, "\n", "")
	cleanURL = strings.ReplaceAll(cleanURL, `"`, "")
	cleanURL = strings.ReplaceAll(cleanURL, `'`, "")
	cleanURL = strings.TrimSuffix(cleanURL, "/")

	// Clean the Service Key too (an invisible \r here breaks the Authorization header!)
	cleanKey := strings.TrimSpace(serviceKey)
	cleanKey = strings.ReplaceAll(cleanKey, "\r", "")
	cleanKey = strings.ReplaceAll(cleanKey, "\n", "")

	// 2. Now safely append the exact storage endpoint
	storageURL := cleanURL + "/storage/v1"

	// 3. Force BOTH required headers for the Supabase API Gateway using the clean key
	customHeaders := map[string]string{
		"apikey":        cleanKey,
		"Authorization": "Bearer " + cleanKey,
	}

	// Initialize the client
	client := storage_go.NewClient(storageURL, cleanKey, customHeaders)

	return &supabaseStorage{
		client: client,
		apiURL: cleanURL, // Keep the clean URL for building the public link later
	}
}

// UploadImage handles streaming the file to Supabase and returning the public URL
func (s *supabaseStorage) UploadImage(bucketName string, filePath string, file io.Reader, contentType string) (string, error) {
	
	// Upload the file to the Supabase bucket
	res, err := s.client.UploadFile(bucketName, filePath, file)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to supabase: %w", err)
	}

	// Log to confirm it worked during development
	log.Printf("File uploaded successfully to Supabase. Response: %v", res)

	// Construct and return the Public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.apiURL, bucketName, filePath)

	return publicURL, nil
}