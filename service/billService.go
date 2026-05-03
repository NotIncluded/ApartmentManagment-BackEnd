package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type BillService interface {
	GenerateMonthlyBill(roomID string, contractID string, recordDate time.Time, dueDate time.Time) (*model.Bill, error)
}

type billService struct {
	billRepo  repository.BillRepository
	roomRepo  *repository.RoomRepository        
	usageRepo *repository.UtilityUsageRepository
	rateRepo  *repository.UtilityRateRepository  
}

func NewBillService(br repository.BillRepository, rr *repository.RoomRepository, ur *repository.UtilityUsageRepository, rate *repository.UtilityRateRepository) BillService {
	return &billService{
		billRepo:  br,
		roomRepo:  rr,
		usageRepo: ur,
		rateRepo:  rate,
	}
}

func (s *billService) GenerateMonthlyBill(roomID string, contractID string, recordDate time.Time, dueDate time.Time) (*model.Bill, error) {
	// 1. Room Expert (BR-02 Validation) - Using your existing FindRoomByID!
	room, err := s.roomRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}
	if room.Status == "Available" {
		return nil, errors.New("BR-02 Violation: Cannot generate bill for an AVAILABLE room")
	}

	// 2. Get Utility Rate - Using your existing FindLatestRate!
	rate, err := s.rateRepo.FindLatestRate()
	if err != nil {
		return nil, errors.New("failed to retrieve active utility rates")
	}

	// 3. Get Utility Usage - Using your existing FindLatestByContract!
	usage, err := s.usageRepo.FindLatestByContract(contractID)
	if err != nil || usage == nil {
		return nil, fmt.Errorf("utility usage data not found for contract %s", contractID)
	}

	// 4. Calculate Usage (BR-12 Validation)
	waterUnits, err := usage.CalculateWaterUsage()
	if err != nil {
		return nil, err
	}
	electricUnits, err := usage.CalculateElectricUsage()
	if err != nil {
		return nil, err
	}

	waterFee := float64(waterUnits) * rate.WaterRate
	electricFee := float64(electricUnits) * rate.ElectricRate

	// 5. Creator: Construct the Bill
	newBill := model.NewBill(
		contractID, // We now pass the real contractID to the bill!
		rate.ID,
		recordDate,
		0, // Rent fee (You can update this later if rent is stored in the contract)
		waterFee,
		electricFee,
		rate.CommonFee,
		dueDate,
	)

	// 6. Save to Database
	if err := s.billRepo.Create(newBill); err != nil {
		return nil, fmt.Errorf("failed to save bill: %w", err)
	}

	return newBill, nil
}