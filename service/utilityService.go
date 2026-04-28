package service

import (
	"errors"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type UtilityService struct {
	rateRepo     *repository.UtilityRateRepository
	usageRepo    *repository.UtilityUsageRepository
	contractRepo *repository.ContractRepository
}

type ConfigureRateRequest struct {
	Period       string  `json:"period"`
	WaterRate    float64 `json:"water_rate"`
	ElectricRate float64 `json:"electric_rate"`
	CommonFee    float64 `json:"common_fee"`
}

type RecordUsageRequest struct {
	ContractID      string    `json:"contract_id"`
	OldWaterUnit    int       `json:"old_water_unit"`
	NewWaterUnit    int       `json:"new_water_unit"`
	OldElectricUnit int       `json:"old_electric_unit"`
	NewElectricUnit int       `json:"new_electric_unit"`
	RecordDate      time.Time `json:"record_date"`
}

func NewUtilityService(
	rateRepo *repository.UtilityRateRepository,
	usageRepo *repository.UtilityUsageRepository,
	contractRepo *repository.ContractRepository,
) *UtilityService {
	return &UtilityService{
		rateRepo:     rateRepo,
		usageRepo:    usageRepo,
		contractRepo: contractRepo,
	}
}

func (s *UtilityService) CreateRate(req ConfigureRateRequest) (*model.UtilityRate, error) {
	if req.Period == "" {
		return nil, errors.New("period is required")
	}
	if req.WaterRate < 0 || req.ElectricRate < 0 || req.CommonFee < 0 {
		return nil, errors.New("rates must be >= 0")
	}

	existing, err := s.rateRepo.FindByPeriod(req.Period)
	if err == nil && existing != nil {
		return nil, errors.New("utility rate for period already exists")
	}

	newRate := model.NewUtilityRate(req.WaterRate, req.ElectricRate, req.CommonFee, req.Period)
	return s.rateRepo.Save(newRate)
}

func (s *UtilityService) ConfigureRate(req ConfigureRateRequest) (*model.UtilityRate, error) {
	if req.Period == "" {
		return nil, errors.New("period is required")
	}
	if req.WaterRate < 0 || req.ElectricRate < 0 || req.CommonFee < 0 {
		return nil, errors.New("rates must be >= 0")
	}

	existing, err := s.rateRepo.FindByPeriod(req.Period)
	if err == nil && existing != nil {
		return nil, errors.New("utility rate for period already exists")
	}

	newRate := model.NewUtilityRate(req.WaterRate, req.ElectricRate, req.CommonFee, req.Period)
	return s.rateRepo.Save(newRate)
}

func (s *UtilityService) GetUtilityRate(period string) (*model.UtilityRate, error) {
	if period != "" {
		return s.rateRepo.FindByPeriod(period)
	}
	return s.rateRepo.FindLatestRate()
}

func (s *UtilityService) GetRateByID(id string) (*model.UtilityRate, error) {
	if id == "" {
		return nil, errors.New("rate id is required")
	}
	return s.rateRepo.FindByID(id)
}

func (s *UtilityService) GetAllRates() ([]model.UtilityRate, error) {
	return s.rateRepo.FindAll()
}

func (s *UtilityService) UpdateRate(id string, req ConfigureRateRequest) (*model.UtilityRate, error) {
	if id == "" {
		return nil, errors.New("rate id is required")
	}

	rate, err := s.rateRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("rate not found")
	}

	if req.WaterRate < 0 || req.ElectricRate < 0 || req.CommonFee < 0 {
		return nil, errors.New("rates must be >= 0")
	}

	// Check if period exists for another rate
	if req.Period != "" && req.Period != rate.Period {
		existing, err := s.rateRepo.FindByPeriod(req.Period)
		if err == nil && existing != nil {
			return nil, errors.New("utility rate for period already exists")
		}
	}

	rate.WaterRate = req.WaterRate
	rate.ElectricRate = req.ElectricRate
	rate.CommonFee = req.CommonFee
	if req.Period != "" {
		rate.Period = req.Period
	}

	return s.rateRepo.Update(rate)
}

func (s *UtilityService) DeleteRate(id string) error {
	if id == "" {
		return errors.New("rate id is required")
	}

	rate, err := s.rateRepo.FindByID(id)
	if err != nil {
		return errors.New("rate not found")
	}

	return s.rateRepo.Delete(rate)
}

func (s *UtilityService) RecordUsage(req RecordUsageRequest) (*model.UtilityUsage, error) {
	if req.ContractID == "" {
		return nil, errors.New("contract id is required")
	}

	_, err := s.contractRepo.FindContractByID(req.ContractID)
	if err != nil {
		return nil, errors.New("contract not found")
	}

	if req.NewWaterUnit < 0 || req.OldWaterUnit < 0 {
		return nil, errors.New("water unit must be >= 0")
	}
	if req.NewElectricUnit < 0 || req.OldElectricUnit < 0 {
		return nil, errors.New("electric unit must be >= 0")
	}
	if req.NewWaterUnit < req.OldWaterUnit {
		return nil, errors.New("new water unit must be >= old water unit")
	}
	if req.NewElectricUnit < req.OldElectricUnit {
		return nil, errors.New("new electric unit must be >= old electric unit")
	}

	prev, err := s.usageRepo.FindLatestByContract(req.ContractID)
	if err != nil {
		return nil, err
	}
	if prev != nil {
		if req.OldWaterUnit != prev.NewWaterUnit {
			return nil, errors.New("old water unit does not match previous record's new water unit")
		}
		if req.OldElectricUnit != prev.NewElectricUnit {
			return nil, errors.New("old electric unit does not match previous record's new electric unit")
		}
	}

	if req.RecordDate.IsZero() {
		req.RecordDate = time.Now()
	}

	usage := model.NewUtilityUsage(req.ContractID, req.OldWaterUnit, req.NewWaterUnit, req.OldElectricUnit, req.NewElectricUnit, req.RecordDate)
	return s.usageRepo.Save(usage)
}

func (s *UtilityService) GetUsageByID(id string) (*model.UtilityUsage, error) {
	if id == "" {
		return nil, errors.New("usage id is required")
	}
	usage, err := s.usageRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("usage not found")
	}
	if usage == nil {
		return nil, errors.New("usage not found")
	}
	contract, err := s.contractRepo.FindContractByID(usage.ContractID)
	if err != nil || contract == nil {
		s.usageRepo.Delete(usage)
		return nil, errors.New("usage not found")
	}
	if contract.DeletedAt.Valid {
		s.usageRepo.Delete(usage)
		return nil, errors.New("usage not found")
	}
	return usage, nil
}

func (s *UtilityService) GetUsagesByContract(contractID string) ([]model.UtilityUsage, error) {
	if contractID == "" {
		return nil, errors.New("contract id is required")
	}
	contract, err := s.contractRepo.FindContractByID(contractID)
	if err != nil || contract == nil {
		return nil, errors.New("contract not found")
	}
	if contract.DeletedAt.Valid {
		return nil, errors.New("contract not found")
	}
	return s.usageRepo.FindByContract(contractID)
}

func (s *UtilityService) UpdateUsage(id string, req RecordUsageRequest) (*model.UtilityUsage, error) {
	if id == "" {
		return nil, errors.New("usage id is required")
	}

	usage, err := s.usageRepo.FindByID(id)
	if err != nil || usage == nil {
		return nil, errors.New("usage not found")
	}

	contract, err := s.contractRepo.FindContractByID(usage.ContractID)
	if err != nil || contract == nil {
		s.usageRepo.Delete(usage)
		return nil, errors.New("usage not found")
	}
	if contract.DeletedAt.Valid {
		s.usageRepo.Delete(usage)
		return nil, errors.New("usage not found")
	}

	if req.NewWaterUnit < req.OldWaterUnit {
		return nil, errors.New("new water unit must be >= old water unit")
	}
	if req.NewElectricUnit < req.OldElectricUnit {
		return nil, errors.New("new electric unit must be >= old electric unit")
	}

	prev, err := s.usageRepo.FindLatestByContract(usage.ContractID)
	if err == nil && prev != nil && prev.ID != id {
		if req.OldWaterUnit != prev.NewWaterUnit {
			return nil, errors.New("old water unit does not match previous record's new water unit")
		}
		if req.OldElectricUnit != prev.NewElectricUnit {
			return nil, errors.New("old electric unit does not match previous record's new electric unit")
		}
	}

	usage.OldWaterUnit = req.OldWaterUnit
	usage.NewWaterUnit = req.NewWaterUnit
	usage.OldElectricUnit = req.OldElectricUnit
	usage.NewElectricUnit = req.NewElectricUnit
	if !req.RecordDate.IsZero() {
		usage.RecordDate = req.RecordDate
	}

	return s.usageRepo.Update(usage)
}

func (s *UtilityService) DeleteUsage(id string) error {
	if id == "" {
		return errors.New("usage id is required")
	}

	usage, err := s.usageRepo.FindByID(id)
	if err != nil || usage == nil {
		return errors.New("usage not found")
	}

	return s.usageRepo.Delete(usage)
}

func (s *UtilityService) GetMyUsages(userID string) ([]model.UtilityUsage, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	contract, err := s.contractRepo.FindActiveContractByUserID(userID)
	if err != nil || contract == nil {
		return nil, errors.New("no active contract found")
	}

	return s.usageRepo.FindByContract(contract.ID)
}

func (s *UtilityService) GetMyLatestUsage(userID string) (*model.UtilityUsage, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	contract, err := s.contractRepo.FindActiveContractByUserID(userID)
	if err != nil || contract == nil {
		return nil, errors.New("no active contract found")
	}

	usage, err := s.usageRepo.FindLatestByContract(contract.ID)
	if err != nil {
		return nil, err
	}
	if usage == nil {
		return nil, errors.New("no usage records found")
	}

	return usage, nil
}
