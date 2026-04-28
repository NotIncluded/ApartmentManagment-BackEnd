package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UtilityUsage struct {
	ID              string         `gorm:"type:char(36);primaryKey" json:"usage_id"`
	ContractID      string         `json:"contract_id" gorm:"not null;index"`
	OldWaterUnit    int            `json:"old_water_unit" gorm:"not null"`
	NewWaterUnit    int            `json:"new_water_unit" gorm:"not null"`
	OldElectricUnit int            `json:"old_electric_unit" gorm:"not null"`
	NewElectricUnit int            `json:"new_electric_unit" gorm:"not null"`
	RecordDate      time.Time      `json:"record_date" gorm:"type:date;not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Contract        *Contract      `gorm:"foreignKey:ContractID;references:ID" json:"-"`
}

func (UtilityUsage) TableName() string {
	return "utility_usages"
}

func (uu *UtilityUsage) BeforeCreate(tx *gorm.DB) (err error) {
	uu.ID = uuid.New().String()
	return
}

func NewUtilityUsage(contractID string, oldWater, newWater, oldElectric, newElectric int, recordDate time.Time) *UtilityUsage {
	return &UtilityUsage{
		ContractID:      contractID,
		OldWaterUnit:    oldWater,
		NewWaterUnit:    newWater,
		OldElectricUnit: oldElectric,
		NewElectricUnit: newElectric,
		RecordDate:      recordDate,
	}
}
