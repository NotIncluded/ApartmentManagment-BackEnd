package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/controller"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	roomController *controller.RoomController
	roomServ       *service.RoomService
)

func initRoomController() {
	roomRepo := repository.NewRoomRepository(testDB)
	contractRepo := repository.NewContractRepository(testDB)
	roomServ = service.NewRoomService(roomRepo, contractRepo)
	roomController = controller.NewRoomController(roomServ)
}

func TestRoomController_GetListRoom_STAFF_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	// Create rooms
	room1 := createTestRoom("601", 6, "Available")
	room2 := createTestRoom("602", 6, "Occupied")

	// Create request with STAFF role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Set("user_id", "staff-user-id")

	// Call controller
	roomController.GetListRoom(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), room1.ID)
	assert.Contains(t, w.Body.String(), room2.ID)
	assert.Contains(t, w.Body.String(), "Rooms retrieved successfully")
}

func TestRoomController_GetListRoom_TENANT_Forbidden(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	// Create user and room
	user, err := authService.Register("Tenant Ctrl", "2222222222", "tenantctrl@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("701", 7, "Occupied")

	// Create active contract
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)
	require.NotNil(t, contract)

	// Create request with TENANT role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Set("user_id", user.ID)

	// Call controller
	roomController.GetListRoom(c)

	// Assert - TENANT should be forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

func TestRoomController_GetMyRoom_TENANT_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	// Create user and room
	user, err := authService.Register("Tenant MyRoom", "0000000000", "myroomtenant@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("901", 9, "Occupied")

	// Create active contract
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now.AddDate(0, 1, 0)
	contract, err := createTestContract(user.ID, room.ID, startDate, endDate, "Active")
	require.NoError(t, err)
	require.NotNil(t, contract)

	// Create request with TENANT role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Set("user_id", user.ID)

	// Call controller
	roomController.GetMyRoom(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), room.ID)
	assert.Contains(t, w.Body.String(), "Your room retrieved successfully")
}

func TestRoomController_GetMyRoom_STAFF_Forbidden(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	// Create request with STAFF role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Set("user_id", "staff-user-id")

	// Call controller
	roomController.GetMyRoom(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "you have no active contract")
}

func TestRoomController_GetMyRoom_TENANT_NoActiveContract_Forbidden(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	// Create user with NO contract
	user, err := authService.Register("Tenant NoMyRoom", "9999999999", "nomyroom@test.com", "password123", "TENANT")
	require.NoError(t, err)

	// Create request with TENANT role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Set("user_id", user.ID)

	// Call controller
	roomController.GetMyRoom(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied: you have no active contract")
}
