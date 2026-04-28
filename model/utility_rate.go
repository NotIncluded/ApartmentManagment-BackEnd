package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UtilityRate struct {
	ID           string    `gorm:"type:char(36);primaryKey" json:"rate_id"`
	WaterRate    float64   `json:"water_rate" gorm:"type:decimal(10,2);not null"`
	ElectricRate float64   `json:"electric_rate" gorm:"type:decimal(10,2);not null"`
	CommonFee    float64   `json:"common_fee" gorm:"type:decimal(10,2);not null"`
	Period       string    `json:"period" gorm:"type:varchar(50)"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Explicit table name (optional but good practice)
func (UtilityRate) TableName() string {
	return "utility_rates"
}

// Auto-generate UUID before insert
func (ur *UtilityRate) BeforeCreate(tx *gorm.DB) (err error) {
	ur.ID = uuid.New().String()
	return
}

// Constructor
func NewUtilityRate(
	waterRate float64,
	electricRate float64,
	commonFee float64,
	period string,
) *UtilityRate {
	return &UtilityRate{
		WaterRate:    waterRate,
		ElectricRate: electricRate,
		CommonFee:    commonFee,
		Period:       period,
	}
}
