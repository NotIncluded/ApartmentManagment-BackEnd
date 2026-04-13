package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomStatus string

const (
	StatusAvailable   RoomStatus = "Available"
	StatusOccupied    RoomStatus = "Occupied"
	StatusMaintenance RoomStatus = "Maintenance"
)

type Room struct {
	ID         string         `gorm:"type:char(36);primaryKey" json:"room_id"`
	RoomNumber string         `json:"room_number" gorm:"not null"`
	Level      int            `json:"level" gorm:"not null"`
	Status     string         `json:"status" gorm:"not null;check:status IN ('Available','Occupied','Maintenance')"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Room) TableName() string {
	return "rooms"
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}

func NewRoom(roomNumber string, level int, status string) *Room {
	return &Room{
		RoomNumber: roomNumber,
		Level:      level,
		Status:     status,
	}
}
