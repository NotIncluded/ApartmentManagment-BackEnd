package service

import (
	"errors"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type RoomService struct {
	roomRepo     *repository.RoomRepository
	contractRepo *repository.ContractRepository
}

func NewRoomService(roomRepo *repository.RoomRepository, contractRepo *repository.ContractRepository) *RoomService {
	return &RoomService{
		roomRepo:     roomRepo,
		contractRepo: contractRepo,
	}
}

// GetListRoom returns all rooms (for STAFF only)
func (s *RoomService) GetListRoom() ([]model.Room, error) {
	rooms, err := s.roomRepo.FindAllRooms()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetRoomByUserID returns the room for a tenant via their active contract
// Returns error if no active contract exists
func (s *RoomService) GetRoomByUserID(userID string) (*model.Room, error) {
	// Find active contract for the tenant
	contract, err := s.contractRepo.FindActiveContractByUserID(userID)
	if err != nil {
		return nil, errors.New("tenant has no active contract")
	}

	// Get the room using the contract's RoomID
	room, err := s.roomRepo.FindRoomByID(contract.RoomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	return room, nil
}
