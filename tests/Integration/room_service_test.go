package integration

import (
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	roomRepo     *repository.RoomRepository
	contractRepo *repository.ContractRepository
	roomService  *service.RoomService
)

func initRoomServices() {
	roomRepo = repository.NewRoomRepository(testDB)
	contractRepo = repository.NewContractRepository(testDB)
	roomService = service.NewRoomService(roomRepo, contractRepo)
}

func setupRoomTestDB() {
	if testDB != nil {
		testDB.AutoMigrate(&model.Room{})
		testDB.AutoMigrate(&model.Contract{})
	}
}

func resetRoomTestDB() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE contracts CASCADE")
		testDB.Exec("TRUNCATE TABLE rooms CASCADE")
	}
}

func cleanupRoomAndContracts() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE contracts CASCADE")
		testDB.Exec("TRUNCATE TABLE rooms CASCADE")
		testDB.Exec("TRUNCATE TABLE users CASCADE")
	}
}

// Helper: Create test room
func createTestRoom(roomNumber string, level int, status string) *model.Room {
	room := model.NewRoom(roomNumber, level, status)
	result := testDB.Create(&room)
	if result.Error != nil {
		panic("Failed to create test room: " + result.Error.Error())
	}
	return room
}

// Helper: Create test contract (returns error if room/user not found due to FK constraint)
func createTestContract(userID, roomID string, startDate, endDate time.Time, status string) (*model.Contract, error) {
	contract := model.NewContract(userID, roomID, startDate, endDate, status)
	contract.Status = status
	result := testDB.Create(&contract)
	if result.Error != nil {
		return nil, result.Error
	}
	return contract, nil
}

func TestRoomService_GetListRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	// Create test rooms
	room1 := createTestRoom("101", 1, "Available")
	room2 := createTestRoom("102", 1, "Occupied")
	room3 := createTestRoom("201", 2, "Maintenance")

	rooms, err := roomService.GetListRoom()

	require.NoError(t, err)
	assert.Equal(t, 3, len(rooms))
	assert.NotNil(t, findRoomByID(rooms, room1.ID))
	assert.NotNil(t, findRoomByID(rooms, room2.ID))
	assert.NotNil(t, findRoomByID(rooms, room3.ID))
}

func TestRoomService_GetListRoom_Empty(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	rooms, err := roomService.GetListRoom()

	require.NoError(t, err)
	assert.Equal(t, 0, len(rooms))
}

func TestRoomService_GetRoomByUserID_WithActiveContract_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	// Create user and room
	user, err := authService.Register("Tenant User", "1234567890", "tenant@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("501", 5, "Occupied")

	// Create active contract
	now := time.Now()
	startDate := now.AddDate(0, -1, 0) // 1 month ago
	endDate := now.AddDate(0, 1, 0)    // 1 month from now
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)
	require.NotNil(t, contract)

	// Test
	retrievedRoom, err := roomService.GetRoomByUserID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, room.ID, retrievedRoom.ID)
	assert.Equal(t, "501", retrievedRoom.RoomNumber)
	assert.Equal(t, 5, retrievedRoom.Level)
}

func TestRoomService_GetRoomByUserID_NoActiveContract_Error(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	// Create user with NO contract
	user, err := authService.Register("Tenant NoContract", "9876543210", "nocontract@test.com", "password123", "TENANT")
	require.NoError(t, err)

	// Test
	_, err = roomService.GetRoomByUserID(user.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active contract found")
}

func TestRoomService_GetRoomByUserID_ExpiredContract_Error(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	// Create user and room
	user, err := authService.Register("Tenant Expired", "5555555555", "expired@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("301", 3, "Available")

	// Create EXPIRED contract
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 1 year ago
	endDate := now.AddDate(-1, 1, 0)   // 11 months ago
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)
	require.NotNil(t, contract)

	// Test
	_, err = roomService.GetRoomByUserID(user.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active contract found")
}

func TestRoomService_GetRoomByUserID_InactiveContract_Error(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	// Create user and room
	user, err := authService.Register("Tenant Inactive", "4444444444", "inactive@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("401", 4, "Available")

	// Create INACTIVE status contract (but within valid date range)
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Pending")
	require.NoError(t, err)
	require.NotNil(t, contract)

	// Test
	_, err = roomService.GetRoomByUserID(user.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active contract found")
}

func TestRoomService_GetRoomByUserID_NonexistentRoom_Error(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	// Create user
	user, err := authService.Register("Tenant NoRoom", "3333333333", "noroom@test.com", "password123", "TENANT")
	require.NoError(t, err)

	// Try to create contract with non-existent room ID
	// This SHOULD fail due to FK constraint - which is the expected behavior
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	_, err = createTestContract(user.ID, "nonexistent-room-id", startDate, endDate, "Active")

	// Assert that contract creation fails (FK constraint)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
}

// ==================== HELPER FUNCTIONS ====================

func findRoomByID(rooms []model.Room, id string) *model.Room {
	for i := range rooms {
		if rooms[i].ID == id {
			return &rooms[i]
		}
	}
	return nil
}
