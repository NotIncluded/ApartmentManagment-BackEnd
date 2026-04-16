package service

import (
	"errors"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type RoomService struct {
	roomRepo     *repository.RoomRepository
	contractRepo *repository.ContractRepository
	userRepo     *repository.UserRepository
}

func NewRoomService(roomRepo *repository.RoomRepository, contractRepo *repository.ContractRepository) *RoomService {
	return &RoomService{
		roomRepo:     roomRepo,
		contractRepo: contractRepo,
	}
}

// Initialize with UserRepository (for AssignRoom and GetRoomTenant)
func (s *RoomService) SetUserRepository(userRepo *repository.UserRepository) {
	s.userRepo = userRepo
}

// CreateRoom creates a new room
func (s *RoomService) CreateRoom(room *model.Room) (*model.Room, error) {
	if room.RoomNumber == "" {
		return nil, errors.New("room number is required")
	}
	if room.Level == 0 {
		return nil, errors.New("level is required")
	}

	err := s.roomRepo.CreateRoom(room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// GetListRoom returns all rooms (for STAFF only)
func (s *RoomService) GetListRoom() ([]model.Room, error) {
	rooms, err := s.roomRepo.FindAllRooms()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// GetRoomByID returns a room by ID
func (s *RoomService) GetRoomByID(roomID string) (*model.Room, error) {
	if roomID == "" {
		return nil, errors.New("room id is required")
	}

	room, err := s.roomRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// UpdateRoom updates room information
func (s *RoomService) UpdateRoom(room *model.Room) (*model.Room, error) {
	if room.ID == "" {
		return nil, errors.New("room id is required")
	}

	err := s.roomRepo.UpdateRoom(room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// DeleteRoom deletes a room - cannot delete if room has any contract
func (s *RoomService) DeleteRoom(roomID string) error {
	if roomID == "" {
		return errors.New("room id is required")
	}

	// Check if room has any active contract
	hasContract, err := s.roomRepo.CheckRoomHasContract(roomID)
	if err != nil {
		return err
	}

	if hasContract {
		return errors.New("cannot delete room with active contract")
	}

	err = s.roomRepo.DeleteRoom(roomID)
	if err != nil {
		return err
	}
	return nil
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

// GetRoomActiveContract returns the active contract of this room
func (s *RoomService) GetRoomActiveContract(roomID string) (*model.Contract, error) {
	if roomID == "" {
		return nil, errors.New("room id is required")
	}

	contract, err := s.contractRepo.FindActiveContractByRoomID(roomID, time.Now())
	if err != nil {
		return nil, err
	}
	return contract, nil
}

// GetRoomContractHistory returns contract history for this room
func (s *RoomService) GetRoomContractHistory(roomID string) ([]model.Contract, error) {
	if roomID == "" {
		return nil, errors.New("room id is required")
	}

	contracts, err := s.contractRepo.FindContractsByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

// GetRoomTenant returns the current tenant via active contract
func (s *RoomService) GetRoomTenant(roomID string) (*model.User, error) {
	if roomID == "" {
		return nil, errors.New("room id is required")
	}

	if s.userRepo == nil {
		return nil, errors.New("user repository not initialized")
	}

	contract, err := s.contractRepo.FindActiveContractByRoomID(roomID, time.Now())
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByID(contract.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AssignRoom assigns a tenant to a room (creates a contract)
// Validations:
// - Room must exist and not be occupied or in maintenance
// - User must exist
func (s *RoomService) AssignRoom(roomID, userID, startDateStr, endDateStr, status string) (*model.Contract, error) {
	// Validate room exists and check status
	room, err := s.roomRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	// Check room status
	if room.Status == string(model.StatusOccupied) {
		return nil, errors.New("room is occupied")
	}
	if room.Status == string(model.StatusMaintenance) {
		return nil, errors.New("room is in maintenance")
	}

	// Validate user exists
	if s.userRepo == nil {
		return nil, errors.New("user repository not initialized")
	}

	_, err = s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start date format (use YYYY-MM-DD)")
	}

	var endDate time.Time
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, errors.New("invalid end date format (use YYYY-MM-DD)")
		}
	}

	// Create contract
	contract := model.NewContract(userID, roomID, startDate, endDate, status)
	err = s.contractRepo.CreateContract(contract)
	if err != nil {
		return nil, err
	}

	return contract, nil
}
