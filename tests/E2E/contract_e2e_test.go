package e2e

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
)

func setupContractTest() {
	setup.ResetTestDB()
}

func cleanupContractTest() {
	setup.ResetTestDB()
}

func registerUserForContractE2E(name, phone, email, role string) (*setup.User, string) {
	user, _ := setup.AuthService.Register(name, phone, email, "password123", role)
	token := GenerateTestToken(user.ID, role)
	return user, token
}

func createTestContractForE2E(userID, roomID string, status string) *setup.Contract {
	now := time.Now()
	contract := &setup.Contract{
		UserID:    userID,
		RoomID:    roomID,
		StartDate: now,
		EndDate:   now.AddDate(0, 6, 0),
		Status:    status,
	}
	setup.TestDB.Create(&contract)
	return contract
}

// ==================== CONTRACT CREATE TESTS ====================

func TestE2E_Contract_Create_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant@test.com", "TENANT")
	room := setup.CreateTestRoom("C101", 1, "Available")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"` + now.Format("2006-01-02") + `","end_date":"` + now.AddDate(0, 6, 0).Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Contract created successfully")
}

func TestE2E_Contract_Create_RoomNotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff2@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant2@test.com", "TENANT")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"nonexistent-room-id","start_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Room not found")
}

func TestE2E_Contract_Create_RoomOccupied(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff3@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant3@test.com", "TENANT")
	room := setup.CreateTestRoom("C102", 1, "Occupied")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Room is not available")
}

func TestE2E_Contract_Create_UserNotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff4@test.com", "STAFF")
	room := setup.CreateTestRoom("C103", 1, "Available")

	now := time.Now()
	data := `{"user_id":"nonexistent-user-id","room_id":"` + room.ID + `","start_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestE2E_Contract_Create_UserHasActiveContract(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff5@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant5@test.com", "TENANT")
	room1 := setup.CreateTestRoom("C104", 1, "Available")
	room2 := setup.CreateTestRoom("C105", 1, "Available")

	now := time.Now()
	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room1.ID,
		StartDate: now,
		EndDate:   now.AddDate(0, 6, 0),
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room2.ID + `","start_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User already has an active contract")
}

func TestE2E_Contract_Create_AfterDeactivatingOld(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff_x@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant_x@test.com", "TENANT")
	room1 := setup.CreateTestRoom("CX01", 1, "Available")
	room2 := setup.CreateTestRoom("CX02", 1, "Available")

	// Create first contract with Active status
	now := time.Now()
	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room1.ID,
		StartDate: now,
		EndDate:   now.AddDate(0, 6, 0),
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)
	room1.Status = "Occupied"
	setup.TestDB.Save(&room1)

	// Step 1: Deactivate old contract
	updateData := `{"status":"Inactive"}`
	w := MakeRequest("PUT", "/contracts/"+contract.ID, staffToken, strings.NewReader(updateData))
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contract updated successfully")

	// Step 2: Create new contract with different room
	createData := `{"user_id":"` + tenant.ID + `","room_id":"` + room2.ID + `","start_date":"` + now.Format("2006-01-02") + `","end_date":"` + now.AddDate(0, 6, 0).Format("2006-01-02") + `","status":"Active"}`
	w = MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(createData))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Contract created successfully")
}

func TestE2E_Contract_Create_InvalidStartDate(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff6@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant6@test.com", "TENANT")
	room := setup.CreateTestRoom("C106", 1, "Available")

	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"invalid-date","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid start date")
}

func TestE2E_Contract_Create_EndDateBeforeStartDate(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff7@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant7@test.com", "TENANT")
	room := setup.CreateTestRoom("C107", 1, "Available")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"` + now.AddDate(0, 6, 0).Format("2006-01-02") + `","end_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "End date must be after start date")
}

func TestE2E_Contract_Create_InactiveStatus(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff8@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant8@test.com", "TENANT")
	room := setup.CreateTestRoom("C108", 1, "Available")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"` + now.Format("2006-01-02") + `","status":"Inactive"}`
	w := MakeRequest("POST", "/contracts/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Contract created successfully")
}

func TestE2E_Contract_Create_Unauthorized(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant9@test.com", "TENANT")
	room := setup.CreateTestRoom("C109", 1, "Available")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", "", strings.NewReader(data))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_Contract_Create_Forbidden_Tenant(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	tenant, token := registerUserForContractE2E("Tenant User", "1111111111", "contract_tenant10@test.com", "TENANT")
	room := setup.CreateTestRoom("C110", 1, "Available")

	now := time.Now()
	data := `{"user_id":"` + tenant.ID + `","room_id":"` + room.ID + `","start_date":"` + now.Format("2006-01-02") + `","status":"Active"}`
	w := MakeRequest("POST", "/contracts/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==================== CONTRACT LIST TESTS ====================

func TestE2E_Contract_GetAll_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff11@test.com", "STAFF")
	tenant1, _ := registerUserForContractE2E("Tenant User 1", "2222222221", "contract_tenant11@test.com", "TENANT")
	tenant2, _ := registerUserForContractE2E("Tenant User 2", "2222222222", "contract_tenant12@test.com", "TENANT")
	room1 := setup.CreateTestRoom("C111", 1, "Available")
	room2 := setup.CreateTestRoom("C112", 1, "Available")

	createTestContractForE2E(tenant1.ID, room1.ID, "Active")
	createTestContractForE2E(tenant2.ID, room2.ID, "Inactive")

	w := MakeRequest("GET", "/contracts/", staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contracts retrieved successfully")
}

func TestE2E_Contract_GetAll_Empty(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff12@test.com", "STAFF")

	w := MakeRequest("GET", "/contracts/", staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contracts retrieved successfully")
}

// ==================== CONTRACT GET BY ID TESTS ====================

func TestE2E_Contract_GetByID_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff13@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant13@test.com", "TENANT")
	room := setup.CreateTestRoom("C113", 1, "Available")

	contract := createTestContractForE2E(tenant.ID, room.ID, "Active")

	w := MakeRequest("GET", "/contracts/"+contract.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contract retrieved successfully")
}

func TestE2E_Contract_GetByID_NotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff14@test.com", "STAFF")

	w := MakeRequest("GET", "/contracts/nonexistent-contract-id", staffToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Contract not found")
}

// ==================== CONTRACT GET BY USER TESTS ====================

func TestE2E_Contract_GetByUser_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff15@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant15@test.com", "TENANT")
	room := setup.CreateTestRoom("C114", 1, "Available")

	createTestContractForE2E(tenant.ID, room.ID, "Active")

	w := MakeRequest("GET", "/contracts/user/"+tenant.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contracts retrieved successfully")
}

func TestE2E_Contract_GetByUser_NotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff16@test.com", "STAFF")

	w := MakeRequest("GET", "/contracts/user/nonexistent-user-id", staffToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

// ==================== CONTRACT GET BY ROOM TESTS ====================

func TestE2E_Contract_GetByRoom_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff17@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant17@test.com", "TENANT")
	room := setup.CreateTestRoom("C115", 1, "Available")

	createTestContractForE2E(tenant.ID, room.ID, "Active")

	w := MakeRequest("GET", "/contracts/room/"+room.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contracts retrieved successfully")
}

func TestE2E_Contract_GetByRoom_NotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff18@test.com", "STAFF")

	w := MakeRequest("GET", "/contracts/room/nonexistent-room-id", staffToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Room not found")
}

// ==================== CONTRACT UPDATE TESTS ====================

func TestE2E_Contract_Update_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff19@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant19@test.com", "TENANT")
	room := setup.CreateTestRoom("C116", 1, "Available")

	contract := createTestContractForE2E(tenant.ID, room.ID, "Active")

	data := `{"status":"Inactive"}`
	w := MakeRequest("PUT", "/contracts/"+contract.ID, staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contract updated successfully")
}

func TestE2E_Contract_Update_ContractNotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff20@test.com", "STAFF")

	data := `{"status":"Inactive"}`
	w := MakeRequest("PUT", "/contracts/nonexistent-contract-id", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Contract not found")
}

func TestE2E_Contract_Update_ChangeToOccupiedRoom(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff21@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant21@test.com", "TENANT")
	room1 := setup.CreateTestRoom("C117", 1, "Available")
	room2 := setup.CreateTestRoom("C118", 1, "Occupied")

	contract := createTestContractForE2E(tenant.ID, room1.ID, "Active")

	data := `{"room_id":"` + room2.ID + `"}`
	w := MakeRequest("PUT", "/contracts/"+contract.ID, staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Target room is not available")
}

func TestE2E_Contract_Update_ChangeToUserWithContract(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff22@test.com", "STAFF")
	tenant1, _ := registerUserForContractE2E("Tenant User 1", "2222222221", "contract_tenant22a@test.com", "TENANT")
	tenant2, _ := registerUserForContractE2E("Tenant User 2", "2222222222", "contract_tenant22b@test.com", "TENANT")
	room1 := setup.CreateTestRoom("C119", 1, "Available")
	room2 := setup.CreateTestRoom("C120", 1, "Available")

	_ = createTestContractForE2E(tenant1.ID, room1.ID, "Active")
	contract2 := createTestContractForE2E(tenant2.ID, room2.ID, "Inactive")

	data := `{"user_id":"` + tenant1.ID + `"}`
	w := MakeRequest("PUT", "/contracts/"+contract2.ID, staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User already has an active contract")
}

func TestE2E_Contract_Update_InvalidDate(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff23@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant23@test.com", "TENANT")
	room := setup.CreateTestRoom("C121", 1, "Available")

	contract := createTestContractForE2E(tenant.ID, room.ID, "Active")

	data := `{"start_date":"invalid-date"}`
	w := MakeRequest("PUT", "/contracts/"+contract.ID, staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid start date")
}

// ==================== CONTRACT DELETE TESTS ====================

func TestE2E_Contract_Delete_Success(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff24@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant24@test.com", "TENANT")
	room := setup.CreateTestRoom("C122", 1, "Available")

	contract := createTestContractForE2E(tenant.ID, room.ID, "Active")

	w := MakeRequest("DELETE", "/contracts/"+contract.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Contract deleted successfully")
}

func TestE2E_Contract_Delete_NotFound(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff25@test.com", "STAFF")

	w := MakeRequest("DELETE", "/contracts/nonexistent-contract-id", staffToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Contract not found")
}

func TestE2E_Contract_Delete_Unauthorized(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant25@test.com", "TENANT")
	room := setup.CreateTestRoom("C123", 1, "Available")

	contract := createTestContractForE2E(tenant.ID, room.ID, "Active")

	w := MakeRequest("DELETE", "/contracts/"+contract.ID, "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_Contract_Delete_Forbidden_Tenant(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	tenant, token := registerUserForContractE2E("Tenant User", "1111111111", "contract_tenant26@test.com", "TENANT")
	room := setup.CreateTestRoom("C124", 1, "Available")

	contract := createTestContractForE2E(tenant.ID, room.ID, "Active")

	w := MakeRequest("DELETE", "/contracts/"+contract.ID, token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==================== CONTRACT EXPIRED TESTS ====================

func TestE2E_Contract_Expired_AutoInactive(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff_exp@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant_exp@test.com", "TENANT")
	room := setup.CreateTestRoom("C200", 1, "Occupied")

	pastDate := time.Now().AddDate(0, -1, 0)
	futureEndDate := time.Now().AddDate(0, -1, 0)

	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: pastDate,
		EndDate:   futureEndDate,
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	w := MakeRequest("GET", "/contracts/", staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), `"status":"Active"`)
}

func TestE2E_Contract_Expired_GetByID(t *testing.T) {
	setupContractTest()
	// defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff_exp2@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant_exp2@test.com", "TENANT")
	room := setup.CreateTestRoom("C201", 1, "Occupied")

	pastDate := time.Now().AddDate(0, -1, 0)
	futureEndDate := time.Now().AddDate(0, -1, 0)

	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: pastDate,
		EndDate:   futureEndDate,
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	w := MakeRequest("GET", "/contracts/"+contract.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), `"status":"Active"`)
}

func TestE2E_Contract_Expired_GetByUser(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff_exp3@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant_exp3@test.com", "TENANT")
	room := setup.CreateTestRoom("C202", 1, "Occupied")

	pastDate := time.Now().AddDate(0, -1, 0)
	futureEndDate := time.Now().AddDate(0, -1, 0)

	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: pastDate,
		EndDate:   futureEndDate,
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)
	room.Status = "Occupied"
	setup.TestDB.Save(&room)

	w := MakeRequest("GET", "/contracts/user/"+tenant.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), `"status":"Active"`)
}

func TestE2E_Contract_Expired_GetByRoom(t *testing.T) {
	setupContractTest()
	defer cleanupContractTest()

	_, staffToken := registerUserForContractE2E("Staff User", "1111111111", "contract_staff_exp4@test.com", "STAFF")
	tenant, _ := registerUserForContractE2E("Tenant User", "2222222222", "contract_tenant_exp4@test.com", "TENANT")
	room := setup.CreateTestRoom("C203", 1, "Occupied")

	pastDate := time.Now().AddDate(0, -1, 0)
	futureEndDate := time.Now().AddDate(0, -1, 0)

	contract := &setup.Contract{
		UserID:    tenant.ID,
		RoomID:    room.ID,
		StartDate: pastDate,
		EndDate:   futureEndDate,
		Status:    "Active",
	}
	setup.TestDB.Create(&contract)

	w := MakeRequest("GET", "/contracts/room/"+room.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, w.Body.String(), `"status":"Active"`)
}
