package repository

import (
	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type BillSlipRepository struct {
	db *gorm.DB
}

func NewBillSlipRepository(db *gorm.DB) *BillSlipRepository {
	return &BillSlipRepository{db: db}
}

func (r *BillSlipRepository) CreateBillSlip(slip *model.BillSlip) error {
	return r.db.Create(slip).Error
}
