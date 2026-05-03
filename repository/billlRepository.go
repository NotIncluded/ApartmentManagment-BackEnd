package repository

import (
	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type BillRepository interface {
	Create(bill *model.Bill) error
	FindByID(id string) (*model.Bill, error)
}

type billRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) BillRepository {
	return &billRepository{db: db}
}

// Create inserts the new generated bill into the database
func (r *billRepository) Create(bill *model.Bill) error {
	return r.db.Create(bill).Error
}

// FindByID will be used later when validating the BillSlip upload!
func (r *billRepository) FindByID(id string) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Where("id = ?", id).First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}