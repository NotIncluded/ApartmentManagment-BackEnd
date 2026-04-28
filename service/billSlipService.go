package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/PunMung-66/ApartmentSys/internal/storage"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type BillSlipService struct {
	repo          *repository.BillSlipRepository
	storageClient storage.StorageService // Replaced MinioClient with our new interface
}

// Updated constructor to accept the new storage client
func NewBillSlipService(repo *repository.BillSlipRepository, storageClient storage.StorageService) *BillSlipService {
	return &BillSlipService{
		repo:          repo,
		storageClient: storageClient,
	}
}

// Updated parameters to match the Controller (file stream and fileName)
func (s *BillSlipService) UploadSlip(ctx context.Context, billID, roomID string, file io.Reader, fileName, contentType string) (string, error) {
	// 1. Generate a unique object path for Supabase
	// We extract the original extension (e.g., .jpg, .png) and append a UNIX timestamp.
	// This prevents browser caching issues and stops files from overwriting each other.
	ext := filepath.Ext(fileName)
	objectPath := fmt.Sprintf("slips/%s-%s-%d%s", billID, roomID, time.Now().Unix(), ext)

	// 2. Upload to Supabase Storage
	// IMPORTANT: You must create a bucket named "apartment-assets" in your Supabase Dashboard
	bucketName := "apartment-assets" 
	slipURL, err := s.storageClient.UploadImage(bucketName, objectPath, file, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to storage: %w", err)
	}

	// 3. Save the resulting public URL to your AWS PostgreSQL Database
	slip := &model.BillSlip{
		ID:      billID + "-" + roomID + "-slip", // Keeping your original ID logic
		BillID:  billID,
		RoomID:  roomID,
		SlipURL: slipURL, // Storing the Supabase URL we just generated
	}

	err = s.repo.CreateBillSlip(slip)
	if err != nil {
		return "", fmt.Errorf("failed to save slip record to database: %w", err)
	}

	return slipURL, nil
}