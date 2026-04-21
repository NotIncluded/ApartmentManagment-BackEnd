package service

import (
	"context"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/internal/minio"
)

type BillSlipService struct {
	repo        *repository.BillSlipRepository
	minioClient *minio.MinioClient
}

func NewBillSlipService(repo *repository.BillSlipRepository, minioClient *minio.MinioClient) *BillSlipService {
	return &BillSlipService{repo: repo, minioClient: minioClient}
}

func (s *BillSlipService) UploadSlip(ctx context.Context, billID, roomID, filePath, contentType string) (string, error) {
	objectName := billID + "-" + roomID + "-slip"
	slipURL, err := s.minioClient.UploadFile(ctx, objectName, filePath, contentType)
	if err != nil {
		return "", err
	}
	slip := &model.BillSlip{
		ID:      objectName,
		BillID:  billID,
		RoomID:  roomID,
		SlipURL: slipURL,
	}
	err = s.repo.CreateBillSlip(slip)
	if err != nil {
		return "", err
	}
	return slipURL, nil
}
