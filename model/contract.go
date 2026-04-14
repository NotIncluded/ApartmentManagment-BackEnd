package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContractStatus string

const (
	ContractStatusInactive ContractStatus = "Inactive"
	ContractStatusActive   ContractStatus = "Active"
)

type Contract struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"contract_id"`
	UserID    string         `json:"user_id" gorm:"not null"`
	RoomID    string         `json:"room_id" gorm:"not null"`
	StartDate time.Time      `json:"start_date" gorm:"type:date;not null"`
	EndDate   time.Time      `json:"end_date" gorm:"type:date"`
	Status    string         `json:"status" gorm:"not null;check:status IN ('Inactive','Active')"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	// Relations
	User *User `gorm:"foreignKey:UserID" json:"-"`
	Room *Room `gorm:"foreignKey:RoomID" json:"-"`
}

func (Contract) TableName() string {
	return "contracts"
}

func (c *Contract) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}

func NewContract(userID, roomID string, startDate, endDate time.Time, status string) *Contract {
	return &Contract{
		UserID:    userID,
		RoomID:    roomID,
		StartDate: startDate,
		EndDate:   endDate,
		Status:    status,
	}
}
