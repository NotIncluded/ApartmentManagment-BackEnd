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
	
	// ✅ FIX 1: Added default:'Unpaid' to satisfy BR-07
	Status         string         `json:"status" gorm:"not null;default:'Unpaid';check:status IN ('Unpaid','WaitingApproval','Paid','Rejected')"`
	
	DueDate        time.Time      `json:"due_date" gorm:"type:date;not null"`
	
	// ✅ FIX 4: Removed redundant CreatedDate. GORM uses CreatedAt automatically!
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	
	// Relations
	Contract       *Contract      `gorm:"foreignKey:ContractID" json:"-"`
	UtilityRate    *UtilityRate   `gorm:"foreignKey:RateID" json:"-"`
	
	// ✅ FIX 3: Link to the upload feature we just built!
	BillSlip       *BillSlip      `gorm:"foreignKey:BillID" json:"bill_slip,omitempty"` 
}

func (Bill) TableName() string {
	return "bills"
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}

// ✅ FIX 2: Information Expert method from your Class Diagram
func (b *Bill) CalculateTotal() {
	b.TotalAmount = b.RentFee + b.WaterFee + b.ElectricityFee + b.CommonFee
}

// Updated constructor: Removed totalAmount, status, and createdDate since they are handled internally/by defaults
func NewBill(contractID, rateID string, recordDate time.Time, rentFee, waterFee, electricityFee, commonFee float64, dueDate time.Time) *Bill {
	bill := &Bill{
		ContractID:     contractID,
		RateID:         rateID,
		RecordDate:     recordDate,
		RentFee:        rentFee,
		WaterFee:       waterFee,
		ElectricityFee: electricityFee,
		CommonFee:      commonFee,
		DueDate:        dueDate,
		// Status is omitted because GORM will default it to "Unpaid"
	}
	
	// Calculate the total immediately upon creation
	bill.CalculateTotal()
	
	return bill
}