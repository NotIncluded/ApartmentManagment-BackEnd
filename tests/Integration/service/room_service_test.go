package service

import (
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	roomRepo         *repository.RoomRepository
	roomContractRepo *repository.ContractRepository
	roomService      *service.RoomService
)

func SetupRoomService() {
	roomRepo = repository.NewRoomRepository(setup.TestDB)
	roomContractRepo = repository.NewContractRepository(setup.TestDB)
	roomService = service.NewRoomService(roomRepo, roomContractRepo)
	setup.RoomService = roomService
}

func initRoomServices() {
	SetupRoomService()
}

func setupRoomTestDB() {
	setup.ResetTestDB()
}

func cleanupRoomAndContracts() {
	setup.ResetTestDB()
}

func createTestRoom(roomNumber string, level int, status string) *model.Room {
	return setup.CreateTestRoom(roomNumber, level, status)
}

func createTestContract(userID, roomID string, startDate, endDate time.Time, status string) (*model.Contract, error) {
	contract := model.NewContract(userID, roomID, startDate, endDate, status)
	contract.Status = status
	result := setup.TestDB.Create(&contract)
	if result.Error != nil {
		return nil, result.Error
	}
	return contract, nil
}

func findRoomByID(rooms []model.Room, id string) *model.Room {
	for i := range rooms {
		if rooms[i].ID == id {
			return &rooms[i]
		}
	}
	return nil
}

// ==================== ORIGINAL TESTS ====================

func TestRoomService_GetListRoom_Success(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

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
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	rooms, err := roomService.GetListRoom()

	require.NoError(t, err)
	assert.Equal(t, 0, len(rooms))
}

func TestRoomService_GetRoomByUserID_WithActiveContract_Success(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant User", "1234567890", "tenant@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("501", 5, "Occupied")
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)
	require.NotNil(t, contract)

	retrievedRoom, err := roomService.GetRoomByUserID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, room.ID, retrievedRoom.ID)
	assert.Equal(t, "501", retrievedRoom.RoomNumber)
	assert.Equal(t, 5, retrievedRoom.Level)
}

func TestRoomService_GetRoomByUserID_NoActiveContract_Error(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant NoContract", "9876543210", "nocontract@test.com", "password123", "TENANT")
	require.NoError(t, err)

	_, err = roomService.GetRoomByUserID(user.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant has no active contract")
}

func TestRoomService_GetRoomByUserID_ExpiredContract_Error(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant Expired", "5555555555", "expired@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("301", 3, "Available")
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0)
	endDate := now.AddDate(-1, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)
	require.NotNil(t, contract)

	_, err = roomService.GetRoomByUserID(user.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant has no active contract")
}

func TestRoomService_GetRoomByUserID_InactiveContract_Error(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant Inactive", "4444444444", "inactive@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("401", 4, "Available")
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, string(model.ContractStatusInactive))
	require.NoError(t, err)
	require.NotNil(t, contract)

	_, err = roomService.GetRoomByUserID(user.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tenant has no active contract")
}

func TestRoomService_GetRoomByUserID_NonexistentRoom_Error(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant NoRoom", "3333333333", "noroom@test.com", "password123", "TENANT")
	require.NoError(t, err)

	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	_, err = createTestContract(user.ID, "nonexistent-room-id", startDate, endDate, "Active")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
}

// ==================== CREATE ROOM ====================

func TestRoomService_CreateRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	newRoom := model.NewRoom("1001", 10, "Available")
	createdRoom, err := roomService.CreateRoom(newRoom)

	require.NoError(t, err)
	assert.Equal(t, "1001", createdRoom.RoomNumber)
	assert.Equal(t, 10, createdRoom.Level)
}

func TestRoomService_CreateRoom_Validation_Error(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	_, err := roomService.CreateRoom(model.NewRoom("", 10, "Available"))
	assert.Error(t, err)

	_, err = roomService.CreateRoom(model.NewRoom("1001", 0, "Available"))
	assert.Error(t, err)
}

// ==================== GET ROOM ====================

func TestRoomService_GetRoomByID_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	room := createTestRoom("1101", 11, "Available")
	retrievedRoom, err := roomService.GetRoomByID(room.ID)

	require.NoError(t, err)
	assert.Equal(t, room.ID, retrievedRoom.ID)
}

func TestRoomService_GetRoomByID_NotFound(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	_, err := roomService.GetRoomByID("nonexistent")
	assert.Error(t, err)
}

// ==================== UPDATE ROOM ====================

func TestRoomService_UpdateRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	room := createTestRoom("1201", 12, "Available")
	room.Status = "Occupied"

	updated, err := roomService.UpdateRoom(room)
	require.NoError(t, err)
	assert.Equal(t, "Occupied", updated.Status)
}

// ==================== DELETE ROOM ====================

func TestRoomService_DeleteRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	room := createTestRoom("1301", 13, "Available")
	err := roomService.DeleteRoom(room.ID)
	require.NoError(t, err)
}

func TestRoomService_DeleteRoom_WithContract_Fails(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant", "7777777777", "delete@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("1401", 14, "Occupied")
	now := time.Now()
	_, err = createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")
	require.NoError(t, err)

	err = roomService.DeleteRoom(room.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete room with active contract")
}

// ==================== GET ROOM ACTIVE CONTRACT ====================

func TestRoomService_GetRoomActiveContract_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant", "8888888888", "contract@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("1501", 15, "Occupied")
	now := time.Now()
	contract, err := createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")
	require.NoError(t, err)

	retrieved, err := roomService.GetRoomActiveContract(room.ID)
	require.NoError(t, err)
	assert.Equal(t, contract.ID, retrieved.ID)
}

// ==================== GET ROOM CONTRACT HISTORY ====================

func TestRoomService_GetRoomContractHistory_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant", "6666666666", "history@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("1701", 17, "Occupied")
	now := time.Now()
	createTestContract(user.ID, room.ID, now.AddDate(-1, 0, 0), now.AddDate(-1, 1, 0), "Inactive")
	createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")

	contracts, err := roomService.GetRoomContractHistory(room.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(contracts), 2)
}

// ==================== GET ROOM TENANT ====================

func TestRoomService_GetRoomTenant_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant", "6363636363", "tenant@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("1901", 19, "Occupied")
	now := time.Now()
	createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")

	userRepo := repository.NewUserRepository(setup.TestDB)
	roomService.SetUserRepository(userRepo)

	tenant, err := roomService.GetRoomTenant(room.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, tenant.ID)
}

// ==================== ASSIGN ROOM ====================

func TestRoomService_AssignRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user, err := setup.AuthService.Register("Tenant", "5252525252", "assign@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("2101", 21, "Available")
	userRepo := repository.NewUserRepository(setup.TestDB)
	roomService.SetUserRepository(userRepo)

	now := time.Now()
	contract, err := roomService.AssignRoom(room.ID, user.ID, now.Format("2006-01-02"), now.AddDate(0, 6, 0).Format("2006-01-02"), "Active")

	require.NoError(t, err)
	assert.Equal(t, room.ID, contract.RoomID)
	assert.Equal(t, user.ID, contract.UserID)
}

func TestRoomService_AssignRoom_Errors(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomServices()

	user1, _ := setup.AuthService.Register("Tenant1", "3232323232", "t1@test.com", "password123", "TENANT")
	user2, _ := setup.AuthService.Register("Tenant2", "2121212121", "t2@test.com", "password123", "TENANT")

	userRepo := repository.NewUserRepository(setup.TestDB)
	roomService.SetUserRepository(userRepo)
	now := time.Now()

	// Room not found
	_, err := roomService.AssignRoom("nonexistent", user1.ID, now.Format("2006-01-02"), now.AddDate(0, 6, 0).Format("2006-01-02"), "Active")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "room not found")

	// User not found
	room := createTestRoom("2201", 22, "Available")
	_, err = roomService.AssignRoom(room.ID, "nonexistent", now.Format("2006-01-02"), now.AddDate(0, 6, 0).Format("2006-01-02"), "Active")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// Room occupied
	room2 := createTestRoom("2301", 23, "Occupied")
	createTestContract(user1.ID, room2.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")
	_, err = roomService.AssignRoom(room2.ID, user2.ID, now.Format("2006-01-02"), now.AddDate(0, 6, 0).Format("2006-01-02"), "Active")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "room is occupied")

	// Room in maintenance
	room3 := createTestRoom("2401", 24, "Maintenance")
	_, err = roomService.AssignRoom(room3.ID, user1.ID, now.Format("2006-01-02"), now.AddDate(0, 6, 0).Format("2006-01-02"), "Active")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "room is in maintenance")

	// Invalid date format
	room4 := createTestRoom("2501", 25, "Available")
	_, err = roomService.AssignRoom(room4.ID, user1.ID, "invalid", "2024-12-31", "Active")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start date format")
}
