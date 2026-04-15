package repository

import (
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	contractRepo *repository.ContractRepository
)

func initContractRepository() {
	contractRepo = repository.NewContractRepository(setup.TestDB)
}

func setupContractTestDB() {
	setup.ResetTestDB()
}

func resetContractTestDB() {
	setup.ResetTestDB()
}

func createTestRoomForContract(roomNumber string, level int, status string) *model.Room {
	return setup.CreateTestRoom(roomNumber, level, status)
}

func createTestContractHelper(userID, roomID string, startDate, endDate time.Time, status string) (*model.Contract, error) {
	contract := model.NewContract(userID, roomID, startDate, endDate, status)
	contract.Status = status
	result := setup.TestDB.Create(&contract)
	if result.Error != nil {
		return nil, result.Error
	}
	return contract, nil
}

// ==================== CREATE CONTRACT ====================

func TestContractRepository_CreateContract_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Contract User", "1111111111", "contractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("101", 1, "Available")

	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)
	contract := model.NewContract(user.ID, room.ID, startDate, endDate, "Active")

	err = contractRepo.CreateContract(contract)

	require.NoError(t, err)
	assert.NotEmpty(t, contract.ID)
	assert.Equal(t, user.ID, contract.UserID)
	assert.Equal(t, room.ID, contract.RoomID)
	assert.Equal(t, "Active", contract.Status)
}

func TestContractRepository_CreateContract_WithEndDateNull(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Open Ended", "1212121212", "openended@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("1001", 10, "Occupied")

	startDate := time.Now().AddDate(0, -1, 0)
	var endDate time.Time
	contract := model.NewContract(user.ID, room.ID, startDate, endDate, "Active")
	contract.EndDate = time.Time{}

	err = contractRepo.CreateContract(contract)

	require.NoError(t, err)
	assert.NotEmpty(t, contract.ID)
}

// ==================== FIND ACTIVE CONTRACT BY USER ID ====================

func TestContractRepository_FindActiveContractByUserID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Active Contract User", "2222222222", "activecontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("201", 2, "Occupied")

	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContractHelper(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByUserID(user.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, user.ID, foundContract.UserID)
	assert.Equal(t, "Active", foundContract.Status)
}

func TestContractRepository_FindActiveContractByUserID_NoActiveContract(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("No Contract User", "3333333333", "nocontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByUserID(user.ID)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
	assert.Contains(t, err.Error(), "no active contract found")
}

func TestContractRepository_FindActiveContractByUserID_InactiveContract(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Inactive Contract User", "4444444444", "inactivecontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("301", 3, "Available")

	now := time.Now()
	startDate := now.AddDate(-1, 0, 0)
	endDate := now.AddDate(-11, 0, 0)
	_, err = createTestContractHelper(user.ID, room.ID, startDate, endDate, "Inactive")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByUserID(user.ID)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
	assert.Contains(t, err.Error(), "no active contract found")
}

// ==================== FIND ACTIVE CONTRACT BY ROOM ID ====================

func TestContractRepository_FindActiveContractByRoomID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Room Contract User", "5555555555", "roomcontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("401", 4, "Occupied")

	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContractHelper(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByRoomID(room.ID, now)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, room.ID, foundContract.RoomID)
	assert.Equal(t, "Active", foundContract.Status)
}

func TestContractRepository_FindActiveContractByRoomID_OutsideDateRange(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Out Range User", "6666666666", "outrangeuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("501", 5, "Occupied")

	now := time.Now()
	startDate := now.AddDate(0, 1, 0)
	endDate := now.AddDate(0, 2, 0)
	_, err = createTestContractHelper(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	searchDate := now.AddDate(0, 0, -1)
	foundContract, err := contractRepo.FindActiveContractByRoomID(room.ID, searchDate)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
	assert.Contains(t, err.Error(), "no active contract found for room")
}

// ==================== FIND CONTRACTS BY ROOM ID ====================

func TestContractRepository_FindContractsByRoomID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user1, err := setup.AuthService.Register("User 1", "7777777777", "user1@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user2, err := setup.AuthService.Register("User 2", "8888888888", "user2@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("601", 6, "Available")

	now := time.Now()
	startDate1 := now.AddDate(-2, 0, 0)
	endDate1 := now.AddDate(-1, 0, 0)
	contract1, err := createTestContractHelper(user1.ID, room.ID, startDate1, endDate1, "Inactive")
	require.NoError(t, err)

	startDate2 := now.AddDate(0, -1, 0)
	endDate2 := now.AddDate(0, 1, 0)
	contract2, err := createTestContractHelper(user2.ID, room.ID, startDate2, endDate2, "Active")
	require.NoError(t, err)

	contracts, err := contractRepo.FindContractsByRoomID(room.ID)

	require.NoError(t, err)
	require.Equal(t, 2, len(contracts))
	assert.Contains(t, []string{contract1.ID, contract2.ID}, contracts[0].ID)
	assert.Contains(t, []string{contract1.ID, contract2.ID}, contracts[1].ID)
}

func TestContractRepository_FindContractsByRoomID_NoContracts(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	room := createTestRoomForContract("701", 7, "Available")

	contracts, err := contractRepo.FindContractsByRoomID(room.ID)

	require.NoError(t, err)
	require.Equal(t, 0, len(contracts))
}

// ==================== FIND CONTRACT BY USER ID ====================

func TestContractRepository_FindContractByUserID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Find by User", "9999999999", "findbyuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("801", 8, "Occupied")

	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContractHelper(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindContractByUserID(user.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, user.ID, foundContract.UserID)
}

func TestContractRepository_FindContractByUserID_NotFound(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("No Contracts", "1010101010", "nocontracts@test.com", "password123", "TENANT")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindContractByUserID(user.ID)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
}

// ==================== FIND CONTRACT BY ID ====================

func TestContractRepository_FindContractByID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	user, err := setup.AuthService.Register("Find by ID", "1111111111", "findbyid@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoomForContract("901", 9, "Occupied")

	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContractHelper(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindContractByID(contract.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, user.ID, foundContract.UserID)
	assert.Equal(t, room.ID, foundContract.RoomID)
}

func TestContractRepository_FindContractByID_NotFound(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	foundContract, err := contractRepo.FindContractByID("nonexistent-contract-id")

	assert.Error(t, err)
	assert.Nil(t, foundContract)
}
