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

// 1. The Interface: Tells the Controller exactly what this service can do
type BillSlipService interface {
	UploadSlip(ctx context.Context, billID, roomID string, file io.Reader, fileName, contentType string) (string, error)
}

// 2. The Struct (lowercase 'b'): Holds our database and storage connections
type billSlipService struct {
	repo          *repository.BillSlipRepository
	billRepo      repository.BillRepository      
	storageClient storage.StorageService         
}

// 3. The Constructor: Returns the Interface, builds the Struct
func NewBillSlipService(repo *repository.BillSlipRepository, billRepo repository.BillRepository, storageClient storage.StorageService) BillSlipService {
	return &billSlipService{
		repo:          repo,
		billRepo:      billRepo,               
		storageClient: storageClient,
	}
}

// 4. The Method (lowercase 'b' receiver): The actual logic
func (s *billSlipService) UploadSlip(ctx context.Context, billID, roomID string, file io.Reader, fileName, contentType string) (string, error) {
	
	// PRE-CHECK: Verify the bill actually exists before we waste time uploading an image!
	_, err := s.billRepo.FindByID(billID)
	if err != nil {
		return "", fmt.Errorf("cannot upload slip: bill ID %s does not exist", billID)
	}

	// 1. Generate a unique object path for Supabase
	ext := filepath.Ext(fileName)
	objectPath := fmt.Sprintf("slips/%s-%s-%d%s", billID, roomID, time.Now().Unix(), ext)

	// 2. Upload to Supabase Storage
	// IMPORTANT: You must create a bucket named "apartment-assets" in your Supabase Dashboard
	bucketName := "apartment-assets" 
	slipURL, err := s.storageClient.UploadImage(bucketName, objectPath, file, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to storage: %w", err)
	}

	// 3. Save the resulting public URL to your PostgreSQL Database
	slip := &model.BillSlip{
		ID:      billID + "-" + roomID + "-slip", 
		BillID:  billID,
		RoomID:  roomID,
		SlipURL: slipURL, 
	}

	err = s.repo.CreateBillSlip(slip)
	if err != nil {
		return "", fmt.Errorf("failed to save slip record to database: %w", err)
	}

	return slipURL, nil
}