package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type BillService interface {
	GenerateMonthlyBill(contractID string, recordDate time.Time, dueDate time.Time) (*model.Bill, error)
}

type billService struct {
	billRepo     repository.BillRepository
	contractRepo repository.ContractRepository       // You will need this
	usageRepo    repository.UtilityUsageRepository   // You will need this
	rateRepo     repository.UtilityRateRepository    // You will need this
}

func NewBillService(br repository.BillRepository, cr repository.ContractRepository, ur repository.UtilityUsageRepository, rr repository.UtilityRateRepository) BillService {
	return &billService{
		billRepo:     br,
		contractRepo: cr,
		usageRepo:    ur,
		rateRepo:     rr,
	}
}

func (s *billService) GenerateMonthlyBill(contractID string, recordDate time.Time, dueDate time.Time) (*model.Bill, error) {
	// 1. Get Contract & Room Data (To check BR-02)
	contract, err := s.contractRepo.FindByID(contractID)
	if err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}
	// Assuming Contract struct has a Room relationship you can check
	if contract.Room.Status == "Available" {
		return nil, errors.New("BR-02 Violation: Cannot generate bill for an AVAILABLE room")
	}

	// 2. Get Utility Rate (The pricing expert)
	// We'll assume you have a function to get the current active rate
	rate, err := s.rateRepo.GetCurrentRate()
	if err != nil {
		return nil, errors.New("failed to retrieve active utility rates")
	}

	// 3. Get Utility Usage (The meter expert)
	usage, err := s.usageRepo.FindByContractAndDate(contractID, recordDate)
	if err != nil {
		return nil, fmt.Errorf("utility usage data not found for contract %s", contractID)
	}

	// 4. Calculate Usage (BR-12 Validation happens inside these methods)
	waterUnits, err := usage.CalculateWaterUsage()
	if err != nil {
		return nil, err // Fails if new unit < old unit
	}
	electricUnits, err := usage.CalculateElectricUsage()
	if err != nil {
		return nil, err
	}

	// Calculate monetary fees
	waterFee := float64(waterUnits) * rate.WaterRate
	electricFee := float64(electricUnits) * rate.ElectricRate

	// 5. Creator: Construct the Bill using your exact constructor
	newBill := model.NewBill(
		contract.ID,
		rate.ID,
		recordDate,
		contract.RentFee, // Assuming rent is tied to the contract
		waterFee,
		electricFee,
		rate.CommonFee,
		dueDate,
	)

	// 6. Save to Database
	if err := s.billRepo.Create(newBill); err != nil {
		return nil, fmt.Errorf("failed to save bill to database: %w", err)
	}

	return newBill, nil
}