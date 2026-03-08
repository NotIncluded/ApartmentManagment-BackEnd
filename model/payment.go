package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "Pending"
	PaymentApproved PaymentStatus = "Approved"
	PaymentRejected PaymentStatus = "Rejected"
)

type Payment struct {
	ID            string         `gorm:"type:char(36);primaryKey" json:"payment_id"`
	BillID        string         `json:"bill_id" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	SlipImagePath string         `json:"slip_image_path"`
	UploadedAt    time.Time      `json:"uploaded_at" gorm:"type:date;not null"`
	ApprovedBy    string         `json:"approved_by"`
	ApprovedAt    time.Time      `json:"approved_at" gorm:"type:date"`
	Status        string         `json:"status" gorm:"not null;check:status IN ('Pending','Approved','Rejected')"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// Relations
	Bill *Bill `gorm:"foreignKey:BillID" json:"-"`
	User *User `gorm:"foreignKey:ApprovedBy" json:"-"`
}

func (Payment) TableName() string {
	return "payments"
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}

func NewPayment(billID string, amount float64, slipImagePath string, uploadedAt time.Time, approvedBy string, approvedAt time.Time, status string) *Payment {
	return &Payment{
		BillID:        billID,
		Amount:        amount,
		SlipImagePath: slipImagePath,
		UploadedAt:    uploadedAt,
		ApprovedBy:    approvedBy,
		ApprovedAt:    approvedAt,
		Status:        status,
	}
}
