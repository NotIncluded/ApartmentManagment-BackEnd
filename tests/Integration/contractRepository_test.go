package integration

import (
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	contractRepo *repository.ContractRepository
)

func initContractRepository() {
	contractRepo = repository.NewContractRepository(testDB)
}

func setupContractTestDB() {
	if testDB != nil {
		testDB.AutoMigrate(&model.User{})
		testDB.AutoMigrate(&model.Room{})
		testDB.AutoMigrate(&model.Contract{})
	}
}

func resetContractTestDB() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE contracts CASCADE")
		testDB.Exec("TRUNCATE TABLE rooms CASCADE")
		testDB.Exec("TRUNCATE TABLE users CASCADE")
	}
}

func cleanupContractTestData(emails []string, roomIDs []string) {
	if testDB != nil {
		for _, email := range emails {
			testDB.Unscoped().Where("email = ?", email).Delete(&model.User{})
		}
		for _, roomID := range roomIDs {
			testDB.Unscoped().Where("id = ?", roomID).Delete(&model.Room{})
		}
		testDB.Exec("TRUNCATE TABLE contracts CASCADE")
	}
}

// TestContractRepository_CreateContract_Success tests successful contract creation
func TestContractRepository_CreateContract_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Contract User", "1111111111", "contractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("101", 1, "Available")

	// Create contract
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now().AddDate(0, 1, 0)
	contract := model.NewContract(user.ID, room.ID, startDate, endDate, "Active")

	err = contractRepo.CreateContract(contract)

	require.NoError(t, err)
	assert.NotEmpty(t, contract.ID)
	assert.Equal(t, user.ID, contract.UserID)
	assert.Equal(t, room.ID, contract.RoomID)
	assert.Equal(t, "Active", contract.Status)

	defer cleanupContractTestData([]string{"contractuser@test.com"}, []string{room.ID})
}

// TestContractRepository_FindActiveContractByUserID_Success tests finding active contract
func TestContractRepository_FindActiveContractByUserID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Active Contract User", "2222222222", "activecontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("201", 2, "Occupied")

	// Create active contract with dates that cover "now"
	now := time.Now()
	startDate := now.AddDate(0, -1, 0) // 1 month ago
	endDate := now.AddDate(0, 1, 0)    // 1 month from now
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByUserID(user.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, user.ID, foundContract.UserID)
	assert.Equal(t, "Active", foundContract.Status)

	defer cleanupContractTestData([]string{"activecontractuser@test.com"}, []string{room.ID})
}

// TestContractRepository_FindActiveContractByUserID_NoActiveContract tests when no active contract exists
func TestContractRepository_FindActiveContractByUserID_NoActiveContract(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user without any contract
	user, err := authService.Register("No Contract User", "3333333333", "nocontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByUserID(user.ID)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
	assert.Contains(t, err.Error(), "no active contract found")

	defer cleanupContractTestData([]string{"nocontractuser@test.com"}, []string{})
}

// TestContractRepository_FindActiveContractByUserID_InactiveContract tests when only inactive contract exists
func TestContractRepository_FindActiveContractByUserID_InactiveContract(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Inactive Contract User", "4444444444", "inactivecontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("301", 3, "Available")

	// Create inactive contract
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 1 year ago
	endDate := now.AddDate(-11, 0, 0)  // 11 months ago
	_, err = createTestContract(user.ID, room.ID, startDate, endDate, "Inactive")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByUserID(user.ID)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
	assert.Contains(t, err.Error(), "no active contract found")

	defer cleanupContractTestData([]string{"inactivecontractuser@test.com"}, []string{room.ID})
}

// TestContractRepository_FindActiveContractByRoomID_Success tests finding active contract by room
func TestContractRepository_FindActiveContractByRoomID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Room Contract User", "5555555555", "roomcontractuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("401", 4, "Occupied")

	// Create active contract
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindActiveContractByRoomID(room.ID, now)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, room.ID, foundContract.RoomID)
	assert.Equal(t, "Active", foundContract.Status)

	defer cleanupContractTestData([]string{"roomcontractuser@test.com"}, []string{room.ID})
}

// TestContractRepository_FindActiveContractByRoomID_OutsideDateRange tests contract outside date range
func TestContractRepository_FindActiveContractByRoomID_OutsideDateRange(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Out Range User", "6666666666", "outrangeuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("501", 5, "Occupied")

	// Create contract that doesn't cover a specific date
	now := time.Now()
	startDate := now.AddDate(0, 1, 0) // starts in 1 month
	endDate := now.AddDate(0, 2, 0)   // ends in 2 months
	_, err = createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	// Search at a date before contract starts
	searchDate := now.AddDate(0, 0, -1) // 1 day ago
	foundContract, err := contractRepo.FindActiveContractByRoomID(room.ID, searchDate)

	assert.Error(t, err)
	assert.Nil(t, foundContract)
	assert.Contains(t, err.Error(), "no active contract found for room")

	defer cleanupContractTestData([]string{"outrangeuser@test.com"}, []string{room.ID})
}

// TestContractRepository_FindContractsByRoomID_Success tests finding all contracts for a room
func TestContractRepository_FindContractsByRoomID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test users and room
	user1, err := authService.Register("User 1", "7777777777", "user1@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user2, err := authService.Register("User 2", "8888888888", "user2@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("601", 6, "Available")

	// Create multiple contracts for the same room
	now := time.Now()
	startDate1 := now.AddDate(-2, 0, 0)
	endDate1 := now.AddDate(-1, 0, 0)
	contract1, err := createTestContract(user1.ID, room.ID, startDate1, endDate1, "Inactive")
	require.NoError(t, err)

	startDate2 := now.AddDate(0, -1, 0)
	endDate2 := now.AddDate(0, 1, 0)
	contract2, err := createTestContract(user2.ID, room.ID, startDate2, endDate2, "Active")
	require.NoError(t, err)

	contracts, err := contractRepo.FindContractsByRoomID(room.ID)

	require.NoError(t, err)
	require.Equal(t, 2, len(contracts))
	assert.Contains(t, []string{contract1.ID, contract2.ID}, contracts[0].ID)
	assert.Contains(t, []string{contract1.ID, contract2.ID}, contracts[1].ID)

	defer cleanupContractTestData([]string{"user1@test.com", "user2@test.com"}, []string{room.ID})
}

// TestContractRepository_FindContractsByRoomID_NoContracts tests when room has no contracts
func TestContractRepository_FindContractsByRoomID_NoContracts(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	room := createTestRoom("701", 7, "Available")

	contracts, err := contractRepo.FindContractsByRoomID(room.ID)

	require.NoError(t, err)
	require.Equal(t, 0, len(contracts))

	defer cleanupContractTestData([]string{}, []string{room.ID})
}

// TestContractRepository_FindContractByUserID_Success tests finding any contract for user
func TestContractRepository_FindContractByUserID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Find by User", "9999999999", "findbyuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("801", 8, "Occupied")

	// Create contract
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindContractByUserID(user.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, user.ID, foundContract.UserID)

	defer cleanupContractTestData([]string{"findbyuser@test.com"}, []string{room.ID})
}

// TestContractRepository_FindContractByUserID_NotFound tests when user has no contracts
func TestContractRepository_FindContractByUserID_NotFound(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create user without contracts
	user, err := authService.Register("No Contracts", "1010101010", "nocontracts@test.com", "password123", "TENANT")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindContractByUserID(user.ID)

	assert.Error(t, err)
	assert.Nil(t, foundContract)

	defer cleanupContractTestData([]string{"nocontracts@test.com"}, []string{})
}

// TestContractRepository_FindContractByID_Success tests finding contract by ID
func TestContractRepository_FindContractByID_Success(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Find by ID", "1111111111", "findbyid@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("901", 9, "Occupied")

	// Create contract
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)

	foundContract, err := contractRepo.FindContractByID(contract.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundContract)
	assert.Equal(t, contract.ID, foundContract.ID)
	assert.Equal(t, user.ID, foundContract.UserID)
	assert.Equal(t, room.ID, foundContract.RoomID)

	defer cleanupContractTestData([]string{"findbyid@test.com"}, []string{room.ID})
}

// TestContractRepository_FindContractByID_NotFound tests when contract ID doesn't exist
func TestContractRepository_FindContractByID_NotFound(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	foundContract, err := contractRepo.FindContractByID("nonexistent-contract-id")

	assert.Error(t, err)
	assert.Nil(t, foundContract)
}

// TestContractRepository_CreateContract_WithEndDateNull tests contract with null end date (open-ended)
func TestContractRepository_CreateContract_WithEndDateNull(t *testing.T) {
	setupContractTestDB()
	defer resetContractTestDB()
	initContractRepository()

	// Create test user and room
	user, err := authService.Register("Open Ended", "1212121212", "openended@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("1001", 10, "Occupied")

	// Create contract with null end date
	startDate := time.Now().AddDate(0, -1, 0)
	var endDate time.Time // zero value for null
	contract := model.NewContract(user.ID, room.ID, startDate, endDate, "Active")
	contract.EndDate = time.Time{} // explicitly set to zero value

	err = contractRepo.CreateContract(contract)

	require.NoError(t, err)
	assert.NotEmpty(t, contract.ID)

	defer cleanupContractTestData([]string{"openended@test.com"}, []string{room.ID})
}
