package service

import (
	"errors"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type ContractService struct {
	contractRepo *repository.ContractRepository
	roomRepo     *repository.RoomRepository
	userRepo     *repository.UserRepository
}

func NewContractService(contractRepo *repository.ContractRepository, roomRepo *repository.RoomRepository) *ContractService {
	return &ContractService{
		contractRepo: contractRepo,
		roomRepo:     roomRepo,
	}
}

func (s *ContractService) SetUserRepository(userRepo *repository.UserRepository) {
	s.userRepo = userRepo
}

func (s *ContractService) CreateContract(userID, roomID, startDateStr, endDateStr, status string) (*model.Contract, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if roomID == "" {
		return nil, errors.New("room id is required")
	}
	if status == "" {
		return nil, errors.New("status is required")
	}
	if status != "Active" && status != "Inactive" {
		return nil, errors.New("status must be Active or Inactive")
	}

	room, err := s.roomRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}
	if room.Status != string(model.StatusAvailable) {
		return nil, errors.New("room is not available")
	}

	if s.userRepo == nil {
		return nil, errors.New("user repository not initialized")
	}

	_, err = s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	activeContract, _ := s.contractRepo.FindActiveContractByUserID(userID)
	if activeContract != nil {
		return nil, errors.New("user already has an active contract")
	}

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
		if endDate.Before(startDate) {
			return nil, errors.New("end date must be after start date")
		}
	}

	contract := model.NewContract(userID, roomID, startDate, endDate, status)
	err = s.contractRepo.CreateContract(contract)
	if err != nil {
		return nil, err
	}

	if status == "Active" {
		room.Status = string(model.StatusOccupied)
		err = s.roomRepo.UpdateRoom(room)
		if err != nil {
			return nil, err
		}
	}

	return contract, nil
}

func (s *ContractService) GetContractByID(contractID string) (*model.Contract, error) {
	if contractID == "" {
		return nil, errors.New("contract id is required")
	}

	contract, err := s.contractRepo.FindContractByID(contractID)
	if err != nil {
		return nil, errors.New("contract not found")
	}

	if contract.Status == "Active" && contract.EndDate.Before(time.Now()) {
		contract.Status = "Inactive"
		s.contractRepo.UpdateContract(contract)

		room, _ := s.roomRepo.FindRoomByID(contract.RoomID)
		if room != nil {
			room.Status = string(model.StatusAvailable)
			s.roomRepo.UpdateRoom(room)
		}
	}

	return contract, nil
}

func (s *ContractService) GetContracts() ([]model.Contract, error) {
	contracts, err := s.contractRepo.FindAllContracts()
	if err != nil {
		return nil, err
	}

	contracts, err = s.processExpiredContracts(contracts)
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

func (s *ContractService) GetContractsByUserID(userID string) ([]model.Contract, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	_, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	contracts, err := s.contractRepo.FindContractsByUserID(userID)
	if err != nil {
		return nil, err
	}

	processedContracts, err := s.processExpiredContracts(contracts)
	if err != nil {
		return nil, err
	}

	return processedContracts, nil
}

func (s *ContractService) GetContractsByRoomID(roomID string) ([]model.Contract, error) {
	if roomID == "" {
		return nil, errors.New("room id is required")
	}

	_, err := s.roomRepo.FindRoomByID(roomID)
	if err != nil {
		return nil, errors.New("room not found")
	}

	contracts, err := s.contractRepo.FindContractsByRoomID(roomID)
	if err != nil {
		return nil, err
	}

	processedContracts, err := s.processExpiredContracts(contracts)
	if err != nil {
		return nil, err
	}

	return processedContracts, nil
}

func (s *ContractService) GetActiveContractByUserID(userID string) (*model.Contract, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	contract, err := s.contractRepo.FindActiveContractByUserID(userID)
	if err != nil {
		return nil, errors.New("no active contract found")
	}

	return contract, nil
}

func (s *ContractService) GetActiveContractByRoomID(roomID string) (*model.Contract, error) {
	if roomID == "" {
		return nil, errors.New("room id is required")
	}

	contract, err := s.contractRepo.FindActiveContractByRoomID(roomID, time.Now())
	if err != nil {
		return nil, errors.New("no active contract found")
	}

	return contract, nil
}

func (s *ContractService) UpdateContract(contractID, userID, roomID, startDateStr, endDateStr, status string) (*model.Contract, error) {
	if contractID == "" {
		return nil, errors.New("contract id is required")
	}

	contract, err := s.contractRepo.FindContractByID(contractID)
	if err != nil {
		return nil, errors.New("contract not found")
	}

	oldRoomID := contract.RoomID
	oldUserID := contract.UserID
	oldStatus := contract.Status

	if userID != "" && userID != oldUserID {
		_, err = s.userRepo.FindUserByID(userID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		activeContract, _ := s.contractRepo.FindActiveContractByUserID(userID)
		if activeContract != nil && activeContract.ID != contractID {
			return nil, errors.New("user already has an active contract")
		}

		contract.UserID = userID
	}

	if roomID != "" && roomID != oldRoomID {
		room, err := s.roomRepo.FindRoomByID(roomID)
		if err != nil {
			return nil, errors.New("room not found")
		}
		if room.Status != string(model.StatusAvailable) {
			return nil, errors.New("target room is not available")
		}

		contract.RoomID = roomID
	}

	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, errors.New("invalid start date format (use YYYY-MM-DD)")
		}
		contract.StartDate = startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, errors.New("invalid end date format (use YYYY-MM-DD)")
		}
		if endDate.Before(contract.StartDate) {
			return nil, errors.New("end date must be after start date")
		}
		contract.EndDate = endDate
	}

	if status != "" {
		if status != "Active" && status != "Inactive" {
			return nil, errors.New("status must be Active or Inactive")
		}
		contract.Status = status
	}

	err = s.contractRepo.UpdateContract(contract)
	if err != nil {
		return nil, err
	}

	if oldStatus == "Active" && status == "Inactive" {
		oldRoom, _ := s.roomRepo.FindRoomByID(oldRoomID)
		if oldRoom != nil {
			oldRoom.Status = string(model.StatusAvailable)
			s.roomRepo.UpdateRoom(oldRoom)
		}
	}

	if status == "Active" || (oldStatus == "Inactive" && contract.Status == "Active") {
		if roomID != "" && roomID != oldRoomID {
			newRoom, _ := s.roomRepo.FindRoomByID(roomID)
			if newRoom != nil {
				newRoom.Status = string(model.StatusOccupied)
				s.roomRepo.UpdateRoom(newRoom)
			}
		}
		if contract.Status == "Active" {
			currentRoom, _ := s.roomRepo.FindRoomByID(contract.RoomID)
			if currentRoom != nil {
				currentRoom.Status = string(model.StatusOccupied)
				s.roomRepo.UpdateRoom(currentRoom)
			}
		}
	}

	return contract, nil
}

func (s *ContractService) DeleteContract(contractID string) error {
	if contractID == "" {
		return errors.New("contract id is required")
	}

	contract, err := s.contractRepo.FindContractByID(contractID)
	if err != nil {
		return errors.New("contract not found")
	}

	roomID := contract.RoomID

	err = s.contractRepo.DeleteContract(contractID)
	if err != nil {
		return err
	}

	room, _ := s.roomRepo.FindRoomByID(roomID)
	if room != nil {
		room.Status = string(model.StatusAvailable)
		s.roomRepo.UpdateRoom(room)
	}

	return nil
}

func (s *ContractService) CheckExpiredContracts() error {
	contracts, err := s.contractRepo.FindAllContracts()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, contract := range contracts {
		if contract.Status == "Active" && contract.EndDate.Before(now) {
			contract.Status = "Inactive"
			s.contractRepo.UpdateContract(&contract)

			room, _ := s.roomRepo.FindRoomByID(contract.RoomID)
			if room != nil {
				room.Status = string(model.StatusAvailable)
				s.roomRepo.UpdateRoom(room)
			}
		}
	}

	return nil
}

func (s *ContractService) HandleUserDeletion(userID string) error {
	if userID == "" {
		return errors.New("user id is required")
	}

	contracts, err := s.contractRepo.FindContractsByUserID(userID)
	if err != nil {
		return err
	}

	for _, contract := range contracts {
		if contract.Status == "Active" {
			contract.Status = "Inactive"
			s.contractRepo.UpdateContract(&contract)

			room, _ := s.roomRepo.FindRoomByID(contract.RoomID)
			if room != nil {
				room.Status = string(model.StatusAvailable)
				s.roomRepo.UpdateRoom(room)
			}
		}
	}

	return nil
}

func (s *ContractService) HandleRoomDeletion(roomID string) error {
	if roomID == "" {
		return errors.New("room id is required")
	}

	contracts, err := s.contractRepo.FindContractsByRoomID(roomID)
	if err != nil {
		return err
	}

	for _, contract := range contracts {
		if contract.Status == "Active" {
			contract.Status = "Inactive"
			s.contractRepo.UpdateContract(&contract)
		}
	}

	return nil
}

func (s *ContractService) processExpiredContracts(contracts []model.Contract) ([]model.Contract, error) {
	now := time.Now()
	updated := false

	for i := range contracts {
		contract := &contracts[i]
		if contract.Status == "Active" && !contract.EndDate.IsZero() && contract.EndDate.Before(now) {
			contract.Status = "Inactive"
			s.contractRepo.UpdateContract(contract)

			room, _ := s.roomRepo.FindRoomByID(contract.RoomID)
			if room != nil {
				room.Status = string(model.StatusAvailable)
				s.roomRepo.UpdateRoom(room)
			}
			updated = true
		}
	}

	if updated {
		return s.contractRepo.FindAllContracts()
	}

	return contracts, nil
}

func (s *ContractService) CheckAndProcessExpired() error {
	return s.CheckExpiredContracts()
}
