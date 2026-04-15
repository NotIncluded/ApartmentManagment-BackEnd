package repository

import (
	"errors"

	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// CreateRoom creates a new room in the database
func (r *RoomRepository) CreateRoom(room *model.Room) error {
	result := r.db.Create(room)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindAllRooms returns all rooms
func (r *RoomRepository) FindAllRooms() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return rooms, nil
}

// FindRoomByID returns a room by ID
func (r *RoomRepository) FindRoomByID(roomID string) (*model.Room, error) {
	var room model.Room
	result := r.db.Where("id = ?", roomID).First(&room)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, result.Error
	}
	return &room, nil
}

// UpdateRoom updates room information
func (r *RoomRepository) UpdateRoom(room *model.Room) error {
	result := r.db.Model(room).Updates(room)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteRoom deletes a room by ID
func (r *RoomRepository) DeleteRoom(roomID string) error {
	result := r.db.Where("id = ?", roomID).Delete(&model.Room{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CheckRoomHasContract checks if room has any contract (active or inactive)
func (r *RoomRepository) CheckRoomHasContract(roomID string) (bool, error) {
	var count int64
	result := r.db.Model(&model.Contract{}).
		Where("room_id = ?", roomID).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
