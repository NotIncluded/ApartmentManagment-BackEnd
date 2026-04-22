package e2e

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
)

func setupRoomTest() {
	setup.ResetTestDB()
}

func cleanupRoomTest() {
	setup.ResetTestDB()
}

func registerUserForE2E(name, phone, email, role string) (*setup.User, string) {
	user, _ := setup.AuthService.Register(name, phone, email, "password123", role)
	token := GenerateTestToken(user.ID, role)
	return user, token
}

// ==================== ROOM CRUD TESTS ====================

func TestE2E_Room_Create_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "1111111111", "e2estaff@test.com", "STAFF")

	data := `{"room_number":"101","level":1,"status":"Available"}`
	w := MakeRequest("POST", "/rooms/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Room created successfully")
	assert.Contains(t, w.Body.String(), "101")
}

func TestE2E_Room_Create_Unauthorized(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	data := `{"room_number":"101","level":1,"status":"Available"}`
	w := MakeRequest("POST", "/rooms/", "", strings.NewReader(data))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_Room_Create_Forbidden_Tenant(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Tenant User", "2222222222", "e2etenant@test.com", "TENANT")

	data := `{"room_number":"101","level":1,"status":"Available"}`
	w := MakeRequest("POST", "/rooms/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_Room_GetList_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "3333333333", "e2estaff2@test.com", "STAFF")

	setup.CreateTestRoom("201", 2, "Available")
	setup.CreateTestRoom("202", 2, "Occupied")

	w := MakeRequest("GET", "/rooms/", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Rooms retrieved successfully")
	assert.Contains(t, w.Body.String(), "201")
	assert.Contains(t, w.Body.String(), "202")
}

func TestE2E_Room_GetList_Forbidden_Tenant(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Tenant User", "4444444444", "e2etenant2@test.com", "TENANT")

	w := MakeRequest("GET", "/rooms/", token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_Room_GetByID_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "5555555555", "e2estaff3@test.com", "STAFF")

	room := setup.CreateTestRoom("301", 3, "Available")

	w := MakeRequest("GET", "/rooms/"+room.ID, token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Room retrieved successfully")
	assert.Contains(t, w.Body.String(), "301")
}

func TestE2E_Room_GetByID_NotFound(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "6666666666", "e2estaff4@test.com", "STAFF")

	w := MakeRequest("GET", "/rooms/nonexistent-id", token, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Room not found")
}

func TestE2E_Room_Update_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "7777777777", "e2estaff5@test.com", "STAFF")

	room := setup.CreateTestRoom("401", 4, "Available")

	data := `{"room_number":"401","level":4,"status":"Occupied"}`
	w := MakeRequest("PUT", "/rooms/"+room.ID, token, strings.NewReader(data))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Room updated successfully")
}

func TestE2E_Room_Delete_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "8888888888", "e2estaff6@test.com", "STAFF")

	room := setup.CreateTestRoom("501", 5, "Available")

	w := MakeRequest("DELETE", "/rooms/"+room.ID, token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Room deleted successfully")
}

// ==================== ROOM RELATIONSHIP TESTS ====================

func TestE2E_Room_Assign_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, staffToken := registerUserForE2E("Staff User", "9999999991", "e2estaff7@test.com", "STAFF")
	tenant, _ := registerUserForE2E("Tenant User", "1010101010", "e2etenant3@test.com", "TENANT")

	room := setup.CreateTestRoom("601", 6, "Available")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","start_date":"` + now.Format("2006-01-02") + `","end_date":"` + now.AddDate(0, 6, 0).Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/rooms/"+room.ID+"/assign", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Room assigned to tenant successfully")
}

func TestE2E_Room_GetActiveContract_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, staffToken := registerUserForE2E("Staff User", "1111111112", "e2estaff8@test.com", "STAFF")
	tenant, _ := registerUserForE2E("Tenant User", "1212121212", "e2etenant4@test.com", "TENANT")

	room := setup.CreateTestRoom("701", 7, "Occupied")

	now := time.Now()
	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now.AddDate(0, 1, 0),
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	w := MakeRequest("GET", "/rooms/"+room.ID+"/contract", staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Active contract retrieved successfully")
}

func TestE2E_Room_GetContractHistory_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, staffToken := registerUserForE2E("Staff User", "1313131313", "e2estaff9@test.com", "STAFF")
	tenant, _ := registerUserForE2E("Tenant User", "1414141414", "e2etenant5@test.com", "TENANT")

	room := setup.CreateTestRoom("801", 8, "Occupied")

	now := time.Now()
	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now.AddDate(0, 1, 0),
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	w := MakeRequest("GET", "/rooms/"+room.ID+"/contracts", staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contract history retrieved successfully")
}

func TestE2E_Room_GetTenant_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, staffToken := registerUserForE2E("Staff User", "1515151515", "e2estaff10@test.com", "STAFF")
	tenant, _ := registerUserForE2E("Tenant User", "1616161616", "e2etenant6@test.com", "TENANT")

	room := setup.CreateTestRoom("901", 9, "Occupied")

	now := time.Now()
	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now.AddDate(0, 1, 0),
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	w := MakeRequest("GET", "/rooms/"+room.ID+"/tenant", staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Tenant retrieved successfully")
	assert.Contains(t, w.Body.String(), tenant.Name)
}

// ==================== ME (TENANT) ROUTE TESTS ====================

func TestE2E_MyRoom_Tenant_Success(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, staffToken := registerUserForE2E("Staff User", "1717171717", "e2estaff11@test.com", "STAFF")
	tenant, tenantToken := registerUserForE2E("Tenant User", "1818181818", "e2etenant7@test.com", "TENANT")

	room := setup.CreateTestRoom("A01", 10, "Occupied")

	now := time.Now()
	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now.AddDate(0, 1, 0),
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	_ = staffToken // used for room setup
	w := MakeRequest("GET", "/me/room", tenantToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Your room retrieved successfully")
	assert.Contains(t, w.Body.String(), room.ID)
}

func TestE2E_MyRoom_Tenant_NoContract(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Tenant User", "1919191919", "e2etenant8@test.com", "TENANT")

	w := MakeRequest("GET", "/me/room", token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_MyRoom_Forbidden_Staff(t *testing.T) {
	setupRoomTest()
	defer cleanupRoomTest()

	_, token := registerUserForE2E("Staff User", "2020202020", "e2estaff12@test.com", "STAFF")

	w := MakeRequest("GET", "/me/room", token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
