package repository

import (
	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type UtilityRateRepository struct {
	db *gorm.DB
}

type UtilityRateRepositoryInterface interface {
	Save(rate *model.UtilityRate) (*model.UtilityRate, error)
	Update(rate *model.UtilityRate) (*model.UtilityRate, error)
	Delete(rate *model.UtilityRate) error
	FindByID(id string) (*model.UtilityRate, error)
	FindAll() ([]model.UtilityRate, error)
	FindLatestRate() (*model.UtilityRate, error)
	FindByPeriod(period string) (*model.UtilityRate, error)
}

func NewUtilityRateRepository(db *gorm.DB) *UtilityRateRepository {
	return &UtilityRateRepository{db: db}
}

func (t *UtilityRateRepository) Save(rate *model.UtilityRate) (*model.UtilityRate, error) {
	result := t.db.Create(rate)
	return rate, result.Error
}

func (t *UtilityRateRepository) Update(rate *model.UtilityRate) (*model.UtilityRate, error) {
	result := t.db.Save(rate)
	return rate, result.Error
}

func (t *UtilityRateRepository) Delete(rate *model.UtilityRate) error {
	result := t.db.Delete(rate)
	return result.Error
}

func (t *UtilityRateRepository) FindByID(id string) (*model.UtilityRate, error) {
	var rate model.UtilityRate
	result := t.db.Where("id = ?", id).First(&rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rate, nil
}

func (t *UtilityRateRepository) FindAll() ([]model.UtilityRate, error) {
	var rates []model.UtilityRate
	result := t.db.Order("created_at desc").Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}
	return rates, nil
}

func (t *UtilityRateRepository) FindLatestRate() (*model.UtilityRate, error) {
	var rate model.UtilityRate
	// Order by period descending (YYYY-MM lexicographic order works) to get latest
	result := t.db.Order("period desc").First(&rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rate, nil
}

func (t *UtilityRateRepository) FindByPeriod(period string) (*model.UtilityRate, error) {
	var rate model.UtilityRate
	result := t.db.Where("period = ?", period).First(&rate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rate, nil
}
