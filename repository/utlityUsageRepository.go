package repository

import (
	"errors"

	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type UtilityUsageRepository struct {
	db *gorm.DB
}

type UtilityUsageRepositoryInterface interface {
	Save(usage *model.UtilityUsage) (*model.UtilityUsage, error)
	Update(usage *model.UtilityUsage) (*model.UtilityUsage, error)
	Delete(usage *model.UtilityUsage) error
	FindByID(id string) (*model.UtilityUsage, error)
	FindByContract(contractID string) ([]model.UtilityUsage, error)
	FindLatestByContract(contractID string) (*model.UtilityUsage, error)
}

func NewUtilityUsageRepository(db *gorm.DB) *UtilityUsageRepository {
	return &UtilityUsageRepository{db: db}
}

func (r *UtilityUsageRepository) Save(usage *model.UtilityUsage) (*model.UtilityUsage, error) {
	result := r.db.Create(&usage)
	return usage, result.Error
}

func (r *UtilityUsageRepository) Update(usage *model.UtilityUsage) (*model.UtilityUsage, error) {
	result := r.db.Save(&usage)
	return usage, result.Error
}

func (r *UtilityUsageRepository) Delete(usage *model.UtilityUsage) error {
	result := r.db.Delete(usage)
	return result.Error
}

func (r *UtilityUsageRepository) FindByID(id string) (*model.UtilityUsage, error) {
	var usage model.UtilityUsage
	result := r.db.Where("id = ?", id).First(&usage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &usage, nil
}

func (r *UtilityUsageRepository) FindByContract(contractID string) ([]model.UtilityUsage, error) {
	var usages []model.UtilityUsage
	result := r.db.Where("contract_id = ?", contractID).Order("record_date asc").Find(&usages)
	if result.Error != nil {
		return nil, result.Error
	}
	return usages, nil
}

func (r *UtilityUsageRepository) FindLatestByContract(contractID string) (*model.UtilityUsage, error) {
	var usage model.UtilityUsage
	result := r.db.Where("contract_id = ?", contractID).Order("record_date desc").First(&usage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &usage, nil
}
