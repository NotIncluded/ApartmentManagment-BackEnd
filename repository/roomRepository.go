package repository

import (
	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) FindAllRooms() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return rooms, nil
}

func (r *RoomRepository) FindRoomByID(roomID string) (*model.Room, error) {
	var room model.Room
	result := r.db.Where("id = ?", roomID).First(&room)
	if result.Error != nil {
		return nil, result.Error
	}
	return &room, nil
}
