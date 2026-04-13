package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UtilityRate struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"rate_id"`
	WaterRate    float64        `json:"water_rate" gorm:"type:decimal(10,2);not null"`
	ElectricRate float64        `json:"electric_rate" gorm:"type:decimal(10,2);not null"`
	CommonFee    float64        `json:"common_fee" gorm:"type:decimal(10,2);not null"`
	ConfiguredBy string         `json:"configured_by" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// Relation
	User *User `gorm:"foreignKey:ConfiguredBy" json:"-"`
}

func (UtilityRate) TableName() string {
	return "utility_rates"
}

func (ur *UtilityRate) BeforeCreate(tx *gorm.DB) (err error) {
	ur.ID = uuid.New().String()
	return
}

func NewUtilityRate(waterRate, electricRate, commonFee float64, configuredBy string) *UtilityRate {
	return &UtilityRate{
		WaterRate:    waterRate,
		ElectricRate: electricRate,
		CommonFee:    commonFee,
		ConfiguredBy: configuredBy,
	}
}
