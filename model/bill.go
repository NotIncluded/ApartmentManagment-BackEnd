package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillStatus string

const (
	BillUnpaid          BillStatus = "Unpaid"
	BillWaitingApproval BillStatus = "WaitingApproval"
	BillPaid            BillStatus = "Paid"
	BillRejected        BillStatus = "Rejected"
)

type Bill struct {
	ID             string         `gorm:"type:char(36);primaryKey" json:"bill_id"`
	ContractID     string         `json:"contract_id" gorm:"not null"`
	RateID         string         `json:"rate_id" gorm:"not null"`
	RecordDate     time.Time      `json:"record_date" gorm:"type:date;not null"`
	RentFee        float64        `json:"rent_fee" gorm:"type:decimal(10,2);not null"`
	WaterFee       float64        `json:"water_fee" gorm:"type:decimal(10,2);not null"`
	ElectricityFee float64        `json:"electricity_fee" gorm:"type:decimal(10,2);not null"`
	CommonFee      float64        `json:"common_fee" gorm:"type:decimal(10,2);not null"`
	TotalAmount    float64        `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status         string         `json:"status" gorm:"not null;check:status IN ('Unpaid','WaitingApproval','Paid','Rejected')"`
	DueDate        time.Time      `json:"due_date" gorm:"type:date;not null"`
	CreatedDate    time.Time      `json:"created_date" gorm:"type:date;not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// Relations
	Contract    *Contract    `gorm:"foreignKey:ContractID" json:"-"`
	UtilityRate *UtilityRate `gorm:"foreignKey:RateID" json:"-"`
}

func (Bill) TableName() string {
	return "bills"
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

func NewBill(contractID, rateID string, recordDate time.Time, rentFee, waterFee, electricityFee, commonFee, totalAmount float64, status string, dueDate, createdDate time.Time) *Bill {
	return &Bill{
		ContractID:     contractID,
		RateID:         rateID,
		RecordDate:     recordDate,
		RentFee:        rentFee,
		WaterFee:       waterFee,
		ElectricityFee: electricityFee,
		CommonFee:      commonFee,
		TotalAmount:    totalAmount,
		Status:         status,
		DueDate:        dueDate,
		CreatedDate:    createdDate,
	}
}
