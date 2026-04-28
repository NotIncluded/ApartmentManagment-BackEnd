package e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
)

func setupUtilityUsageTest() {
	setup.ResetTestDB()
}

func cleanupUtilityUsageTest() {
	setup.ResetTestDB()
}

func createContractForUtilityUsageE2E(userID, roomID string, startDate, endDate string) (*setup.Contract, error) {
	return setup.CreateTestContract(userID, roomID, startDate, endDate, "Active")
}

func registerUserForUtilityUsageE2E(name, phone, email, role string) (*setup.User, string) {
	user, _ := setup.AuthService.Register(name, phone, email, "password123", role)
	token := GenerateTestToken(user.ID, role)
	return user, token
}

// ==================== UTILITY USAGE CREATE TESTS ====================

func TestE2E_UtilityUsage_Create_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "1111111111", "uu_staff@test.com", "STAFF")
	room := setup.CreateTestRoom("101", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	data := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	w := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Utility usage recorded successfully")
	assert.Contains(t, w.Body.String(), "100")
	assert.Contains(t, w.Body.String(), "500")
}

func TestE2E_UtilityUsage_Create_Unauthorized(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	data := `{"contract_id":"some-id","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	w := MakeRequest("POST", "/utility-usages/", "", strings.NewReader(data))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_UtilityUsage_Create_Forbidden_Tenant(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "2222222222", "uu_tenant@test.com", "TENANT")
	room := setup.CreateTestRoom("102", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	data := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	w := MakeRequest("POST", "/utility-usages/", tenantToken, strings.NewReader(data))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_UtilityUsage_Create_ContractNotFound(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "3333333333", "uu_staff2@test.com", "STAFF")

	data := `{"contract_id":"nonexistent-contract","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	w := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "contract not found")
}

func TestE2E_UtilityUsage_Create_InvalidUnits_NegativeWater(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "4444444444", "uu_staff3@test.com", "STAFF")
	room := setup.CreateTestRoom("103", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	data := `{"contract_id":"` + contract.ID + `","old_water_unit":-10,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	w := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "water unit must be")
}

func TestE2E_UtilityUsage_Create_InvalidUnits_NewLessThanOld(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "5555555555", "uu_staff4@test.com", "STAFF")
	room := setup.CreateTestRoom("104", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	data := `{"contract_id":"` + contract.ID + `","old_water_unit":100,"new_water_unit":50,"old_electric_unit":0,"new_electric_unit":500}`
	w := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "new water unit must be")
}

// ==================== UTILITY USAGE GET BY CONTRACT TESTS ====================

func TestE2E_UtilityUsage_GetByContract_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "6666666666", "uu_staff5@test.com", "STAFF")
	room := setup.CreateTestRoom("105", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))

	w := MakeRequest("GET", "/utility-usages/contract/"+contract.ID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Usages retrieved successfully")
	assert.Contains(t, w.Body.String(), "100")
}

func TestE2E_UtilityUsage_GetByContract_ContractNotFound(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "7777777777", "uu_staff6@test.com", "STAFF")

	w := MakeRequest("GET", "/utility-usages/contract/nonexistent-contract", staffToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ==================== UTILITY USAGE GET BY ID TESTS ====================

func TestE2E_UtilityUsage_GetByID_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "8888888888", "uu_staff7@test.com", "STAFF")
	room := setup.CreateTestRoom("106", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	createResp := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"usage_id":"`)
	if startIdx == -1 {
		t.Fatal("Could not find usage_id in response")
	}
	usageIDStart := startIdx + 12
	endIdx := strings.Index(responseBody[usageIDStart:], `"`)
	usageID := responseBody[usageIDStart : usageIDStart+endIdx]

	w := MakeRequest("GET", "/utility-usages/"+usageID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Usage retrieved successfully")
	assert.Contains(t, w.Body.String(), usageID)
}

func TestE2E_UtilityUsage_GetByID_NotFound(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "9999999999", "uu_staff8@test.com", "STAFF")

	w := MakeRequest("GET", "/utility-usages/nonexistent-id", staffToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestE2E_UtilityUsage_GetByID_Unauthorized(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	w := MakeRequest("GET", "/utility-usages/some-id", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ==================== UTILITY USAGE UPDATE TESTS ====================

func TestE2E_UtilityUsage_Update_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "1010101010", "uu_staff9@test.com", "STAFF")
	room := setup.CreateTestRoom("107", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	createResp := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"usage_id":"`)
	usageIDStart := startIdx + 12
	endIdx := strings.Index(responseBody[usageIDStart:], `"`)
	usageID := responseBody[usageIDStart : usageIDStart+endIdx]

	updateData := `{"old_water_unit":100,"new_water_unit":200,"old_electric_unit":500,"new_electric_unit":1000}`
	w := MakeRequest("PUT", "/utility-usages/"+usageID, staffToken, strings.NewReader(updateData))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Utility usage updated successfully")
	assert.Contains(t, w.Body.String(), "200")
}

func TestE2E_UtilityUsage_Update_NotFound(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "1111101010", "uu_staff10@test.com", "STAFF")

	updateData := `{"old_water_unit":100,"new_water_unit":200,"old_electric_unit":500,"new_electric_unit":1000}`
	w := MakeRequest("PUT", "/utility-usages/nonexistent-id", staffToken, strings.NewReader(updateData))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "usage not found")
}

func TestE2E_UtilityUsage_Update_Unauthorized(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	updateData := `{"old_water_unit":100,"new_water_unit":200,"old_electric_unit":500,"new_electric_unit":1000}`
	w := MakeRequest("PUT", "/utility-usages/some-id", "", strings.NewReader(updateData))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ==================== UTILITY USAGE DELETE TESTS ====================

func TestE2E_UtilityUsage_Delete_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	staff, staffToken := registerUserForUtilityUsageE2E("Staff User", "1212101010", "uu_staff11@test.com", "STAFF")
	room := setup.CreateTestRoom("108", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(staff.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	createResp := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"usage_id":"`)
	usageIDStart := startIdx + 12
	endIdx := strings.Index(responseBody[usageIDStart:], `"`)
	usageID := responseBody[usageIDStart : usageIDStart+endIdx]

	w := MakeRequest("DELETE", "/utility-usages/"+usageID, staffToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Utility usage deleted successfully")
}

func TestE2E_UtilityUsage_Delete_NotFound(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "1313101010", "uu_staff12@test.com", "STAFF")

	w := MakeRequest("DELETE", "/utility-usages/nonexistent-id", staffToken, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "usage not found")
}

func TestE2E_UtilityUsage_Delete_Unauthorized(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	w := MakeRequest("DELETE", "/utility-usages/some-id", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_UtilityUsage_Delete_Forbidden_Tenant(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "1515101010", "uu_staff13@test.com", "STAFF")
	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "1414101010", "uu_tenant2@test.com", "TENANT")
	room := setup.CreateTestRoom("109", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	createResp := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"usage_id":"`)
	usageIDStart := startIdx + 12
	endIdx := strings.Index(responseBody[usageIDStart:], `"`)
	usageID := responseBody[usageIDStart : usageIDStart+endIdx]

	w := MakeRequest("DELETE", "/utility-usages/"+usageID, tenantToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==================== TENANT MY USAGES TESTS ====================

func TestE2E_UtilityUsage_GetMyUsages_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "1616101010", "uu_staff14@test.com", "STAFF")
	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "1717101010", "uu_tenant3@test.com", "TENANT")
	room := setup.CreateTestRoom("110", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))

	w := MakeRequest("GET", "/me/usages", tenantToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Usages retrieved successfully")
}

func TestE2E_UtilityUsage_GetMyUsages_NoContract(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, tenantToken := registerUserForUtilityUsageE2E("Tenant User No Contract", "1818101010", "uu_tenant4@test.com", "TENANT")

	w := MakeRequest("GET", "/me/usages", tenantToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestE2E_UtilityUsage_GetMyUsages_Unauthorized(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	w := MakeRequest("GET", "/me/usages", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ==================== TENANT MY LATEST USAGE TESTS ====================

func TestE2E_UtilityUsage_GetMyLatestUsage_Success(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "1919101010", "uu_staff15@test.com", "STAFF")
	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "2020101010", "uu_tenant5@test.com", "TENANT")
	room := setup.CreateTestRoom("111", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))

	w := MakeRequest("GET", "/me/usages/latest", tenantToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Latest usage retrieved successfully")
	assert.Contains(t, w.Body.String(), "100")
}

func TestE2E_UtilityUsage_GetMyLatestUsage_NoContract(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, tenantToken := registerUserForUtilityUsageE2E("Tenant User No Contract", "2121101010", "uu_tenant6@test.com", "TENANT")

	w := MakeRequest("GET", "/me/usages/latest", tenantToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "no active contract found")
}

func TestE2E_UtilityUsage_GetMyLatestUsage_NoUsage(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "2222101010", "uu_staff16@test.com", "STAFF")
	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "2323101010", "uu_tenant7@test.com", "TENANT")
	room := setup.CreateTestRoom("112", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))
	_ = room

	w := MakeRequest("GET", "/me/usages/latest", tenantToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Latest usage retrieved successfully")
}

func TestE2E_UtilityUsage_GetMyLatestUsage_Unauthorized(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	w := MakeRequest("GET", "/me/usages/latest", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ==================== TENANT FORBIDDEN FROM DIRECT ACCESS TESTS ====================

func TestE2E_UtilityUsage_GetByContract_Forbidden_Tenant(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, _ = registerUserForUtilityUsageE2E("Staff User", "2424101010", "uu_staff17@test.com", "STAFF")
	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "2525101010", "uu_tenant8@test.com", "TENANT")
	room := setup.CreateTestRoom("113", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)
	_ = room

	w := MakeRequest("GET", "/utility-usages/contract/"+contract.ID, tenantToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_UtilityUsage_GetByID_Forbidden_Tenant(t *testing.T) {
	setupUtilityUsageTest()
	defer cleanupUtilityUsageTest()

	_, staffToken := registerUserForUtilityUsageE2E("Staff User", "2626101010", "uu_staff18@test.com", "STAFF")
	tenant, tenantToken := registerUserForUtilityUsageE2E("Tenant User", "2727101010", "uu_tenant9@test.com", "TENANT")
	room := setup.CreateTestRoom("114", 1, "Available")

	contract, err := createContractForUtilityUsageE2E(tenant.ID, room.ID, "2020-01-01", "2030-12-31")
	assert.NoError(t, err)
	_ = room

	createData := `{"contract_id":"` + contract.ID + `","old_water_unit":0,"new_water_unit":100,"old_electric_unit":0,"new_electric_unit":500}`
	createResp := MakeRequest("POST", "/utility-usages/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"usage_id":"`)
	usageIDStart := startIdx + 12
	endIdx := strings.Index(responseBody[usageIDStart:], `"`)
	usageID := responseBody[usageIDStart : usageIDStart+endIdx]

	w := MakeRequest("GET", "/utility-usages/"+usageID, tenantToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
