package repository

import (
	"errors"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type ContractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// CreateContract creates a new contract in the database
func (c *ContractRepository) CreateContract(contract *model.Contract) error {
	result := c.db.Create(contract)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindActiveContractByUserID returns the active contract for a user
// Active contract: Status is "Active" and current date is between StartDate and EndDate
func (c *ContractRepository) FindActiveContractByUserID(userID string) (*model.Contract, error) {
	var contract model.Contract
	now := time.Now()

	result := c.db.Where(
		"user_id = ? AND status = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
		userID,
		"Active",
		now,
		now,
	).First(&contract)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active contract found")
		}
		return nil, result.Error
	}

	return &contract, nil
}

// FindActiveContractByRoomID returns the active contract for a room at given time
func (c *ContractRepository) FindActiveContractByRoomID(roomID string, at time.Time) (*model.Contract, error) {
	var contract model.Contract

	result := c.db.Where(
		"room_id = ? AND status = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
		roomID,
		"Active",
		at,
		at,
	).First(&contract)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active contract found for room")
		}
		return nil, result.Error
	}

	return &contract, nil
}

// FindContractsByRoomID returns all contracts for a room (contract history)
func (c *ContractRepository) FindContractsByRoomID(roomID string) ([]model.Contract, error) {
	var contracts []model.Contract
	result := c.db.Where("room_id = ?", roomID).Find(&contracts)
	if result.Error != nil {
		return nil, result.Error
	}
	return contracts, nil
}

// FindContractByUserID returns any contract for a user (regardless of status)
func (c *ContractRepository) FindContractByUserID(userID string) (*model.Contract, error) {
	var contract model.Contract
	result := c.db.Where("user_id = ?", userID).First(&contract)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contract, nil
}

func (c *ContractRepository) FindContractByID(contractID string) (*model.Contract, error) {
	var contract model.Contract
	result := c.db.Where("id = ?", contractID).First(&contract)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contract, nil
}
