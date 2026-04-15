package integration

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== CREATE ROOM ====================

func TestRoomController_CreateRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	roomData := `{"room_number": "3001", "level": 30, "status": "Available"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Request, _ = http.NewRequest("POST", "/rooms", io.NopCloser(strings.NewReader(roomData)))
	c.Request.Header.Set("Content-Type", "application/json")

	roomController.CreateRoom(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Room created successfully")
}

func TestRoomController_CreateRoom_TenantDenied(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	roomData := `{"room_number": "9999", "level": 99, "status": "Available"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Request, _ = http.NewRequest("POST", "/rooms", io.NopCloser(strings.NewReader(roomData)))
	c.Request.Header.Set("Content-Type", "application/json")

	roomController.CreateRoom(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

// ==================== GET LIST ROOM ====================

func TestRoomController_GetListRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	createTestRoom("3000", 30, "Available")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")

	roomController.GetListRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Rooms retrieved successfully")
}

func TestRoomController_GetListRoom_TenantDenied(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Set("user_id", "some-user-id")

	roomController.GetListRoom(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

// ==================== GET ROOM BY ID ====================

func TestRoomController_GetRoomByID_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	room := createTestRoom("3101", 31, "Available")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.GetRoomByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "3101")
}

func TestRoomController_GetRoomByID_NotFound(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

	roomController.GetRoomByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ==================== UPDATE ROOM ====================

func TestRoomController_UpdateRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	room := createTestRoom("3201", 32, "Available")
	updateData := `{"status": "Occupied"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}
	c.Request, _ = http.NewRequest("PUT", "/rooms/"+room.ID, strings.NewReader(updateData))
	c.Request.Header.Set("Content-Type", "application/json")

	roomController.UpdateRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Occupied")
}

func TestRoomController_UpdateRoom_TenantDenied(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	room := createTestRoom("3211", 32, "Available")
	updateData := `{"status": "Occupied"}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}
	c.Request, _ = http.NewRequest("PUT", "/rooms/"+room.ID, strings.NewReader(updateData))
	c.Request.Header.Set("Content-Type", "application/json")

	roomController.UpdateRoom(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

// ==================== DELETE ROOM ====================

func TestRoomController_DeleteRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	room := createTestRoom("3401", 34, "Available")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.DeleteRoom(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoomController_DeleteRoom_TenantDenied(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	room := createTestRoom("3411", 34, "Available")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.DeleteRoom(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

func TestRoomController_DeleteRoom_WithContract(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	user, err := authService.Register("Tenant", "7777777777", "delete@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("3501", 35, "Occupied")
	now := time.Now()
	createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.DeleteRoom(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Cannot delete room with existing contract")
}

// ==================== GET ROOM ACTIVE CONTRACT ====================

func TestRoomController_GetRoomActiveContract_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	user, err := authService.Register("Tenant", "6666666666", "activecontract@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("3601", 36, "Occupied")
	now := time.Now()
	contract, err := createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.GetRoomActiveContract(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), contract.ID)
}

func TestRoomController_GetRoomActiveContract_TenantDenied(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	room := createTestRoom("3611", 36, "Available")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.GetRoomActiveContract(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

// ==================== GET ROOM CONTRACT HISTORY ====================

func TestRoomController_GetRoomContractHistory_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	user, err := authService.Register("Tenant", "5555555555", "history@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("3801", 38, "Occupied")
	now := time.Now()
	createTestContract(user.ID, room.ID, now.AddDate(-1, 0, 0), now.AddDate(-1, 1, 0), "Completed")
	createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.GetRoomContractHistory(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==================== GET ROOM TENANT ====================

func TestRoomController_GetRoomTenant_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()
	log.SetOutput(os.Stdout)

	user, err := authService.Register("Tenant", "4444444444", "gettenant@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("4001", 40, "Occupied")
	now := time.Now()
	createTestContract(user.ID, room.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")

	// Set user repository for GetRoomTenant to work
	userRepo := repository.NewUserRepository(testDB)
	roomServ.SetUserRepository(userRepo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}

	roomController.GetRoomTenant(c)

	assert.Equal(t, http.StatusOK, w.Code)
	t.Log(w.Body.String())
	assert.Contains(t, w.Body.String(), user.ID)
}

// ==================== ASSIGN ROOM ====================

func TestRoomController_AssignRoom_Success(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	user, err := authService.Register("Tenant", "3333333333", "assign@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("4201", 42, "Available")

	data := struct {
		UserID    string `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status"`
	}{UserID: user.ID, StartDate: time.Now().Format("2006-01-02"), EndDate: time.Now().AddDate(0, 6, 0).Format("2006-01-02"), Status: "Active"}

	jsonData, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}
	c.Request, _ = http.NewRequest("POST", "/rooms/"+room.ID+"/assign", strings.NewReader(string(jsonData)))
	c.Request.Header.Set("Content-Type", "application/json")

	roomController.AssignRoom(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Room assigned to tenant successfully")
}

func TestRoomController_AssignRoom_TenantDenied(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	user, err := authService.Register("Tenant", "3333333334", "assigndeny@test.com", "password123", "TENANT")
	require.NoError(t, err)

	room := createTestRoom("4211", 42, "Available")

	data := struct {
		UserID    string `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status"`
	}{UserID: user.ID, StartDate: time.Now().Format("2006-01-02"), EndDate: time.Now().AddDate(0, 6, 0).Format("2006-01-02"), Status: "Active"}

	jsonData, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "TENANT")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}
	c.Request, _ = http.NewRequest("POST", "/rooms/"+room.ID+"/assign", strings.NewReader(string(jsonData)))
	c.Request.Header.Set("Content-Type", "application/json")

	roomController.AssignRoom(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

func TestRoomController_AssignRoom_Errors(t *testing.T) {
	setupRoomTestDB()
	defer cleanupRoomAndContracts()
	initRoomController()

	user1, _ := authService.Register("T1", "1111111111", "t1@test.com", "password123", "TENANT")
	user2, _ := authService.Register("T2", "2222222222", "t2@test.com", "password123", "TENANT")

	now := time.Now()
	startDateStr := now.Format("2006-01-02")
	endDateStr := now.AddDate(0, 6, 0).Format("2006-01-02")

	// Room not found
	data := struct {
		UserID    string `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Status    string `json:"status"`
	}{UserID: user1.ID, StartDate: startDateStr, EndDate: endDateStr, Status: "Active"}
	jsonData, _ := json.Marshal(data)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request, _ = http.NewRequest("POST", "/rooms/nonexistent/assign", strings.NewReader(string(jsonData)))
	c.Request.Header.Set("Content-Type", "application/json")
	roomController.AssignRoom(c)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// User not found
	room := createTestRoom("4301", 43, "Available")
	data.UserID = "nonexistent"
	jsonData, _ = json.Marshal(data)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room.ID}}
	c.Request, _ = http.NewRequest("POST", "/rooms/"+room.ID+"/assign", strings.NewReader(string(jsonData)))
	c.Request.Header.Set("Content-Type", "application/json")
	roomController.AssignRoom(c)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Room occupied
	room2 := createTestRoom("4401", 44, "Occupied")
	createTestContract(user1.ID, room2.ID, now.AddDate(0, -1, 0), now.AddDate(0, 1, 0), "Active")
	data.UserID = user2.ID
	jsonData, _ = json.Marshal(data)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room2.ID}}
	c.Request, _ = http.NewRequest("POST", "/rooms/"+room2.ID+"/assign", strings.NewReader(string(jsonData)))
	c.Request.Header.Set("Content-Type", "application/json")
	roomController.AssignRoom(c)
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "room is occupied")

	// Room in maintenance
	room3 := createTestRoom("4501", 45, "Maintenance")
	data.UserID = user1.ID
	jsonData, _ = json.Marshal(data)
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", "STAFF")
	c.Params = gin.Params{{Key: "id", Value: room3.ID}}
	c.Request, _ = http.NewRequest("POST", "/rooms/"+room3.ID+"/assign", strings.NewReader(string(jsonData)))
	c.Request.Header.Set("Content-Type", "application/json")
	roomController.AssignRoom(c)
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "room is in maintenance")
}
