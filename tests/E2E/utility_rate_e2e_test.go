package e2e

import (
	"net/http"
	"strings"
	"testing"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
)

func setupUtilityRateTest() {
	setup.ResetTestDB()
}

func cleanupUtilityRateTest() {
	setup.ResetTestDB()
}

func registerUserForUtilityRateE2E(name, phone, email, role string) (*setup.User, string) {
	user, _ := setup.AuthService.Register(name, phone, email, "password123", role)
	token := GenerateTestToken(user.ID, role)
	return user, token
}

// ==================== UTILITY RATE CREATE TESTS ====================

func TestE2E_UtilityRate_Create_Success(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1111111111", "ur_staff@test.com", "STAFF")

	data := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-01"}}`
	w := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Utility rate created successfully")
	assert.Contains(t, w.Body.String(), "5.5")
	assert.Contains(t, w.Body.String(), "7.2")
	assert.Contains(t, w.Body.String(), "2.5")
}

func TestE2E_UtilityRate_Create_Unauthorized(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	data := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-01"}}`
	w := MakeRequest("POST", "/utility-rates/", "", strings.NewReader(data))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_UtilityRate_Create_Forbidden_Tenant(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Tenant User", "2222222222", "ur_tenant@test.com", "TENANT")

	data := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-01"}}`
	w := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestE2E_UtilityRate_Create_InvalidRates_NegativeWater(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "3333333333", "ur_staff2@test.com", "STAFF")

	data := `{"water_rate":-5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-01","configured_by":"admin"}`
	w := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid rates")
}

func TestE2E_UtilityRate_Create_DuplicatePeriod(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "4444444444", "ur_staff3@test.com", "STAFF")

	// Create first rate
	data1 := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-02","configured_by":"admin"}`
	w1 := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data1))
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Try to create duplicate period
	data2 := `{"water_rate":6.0,"electric_rate":8.0,"common_fee":3.0,"period":"2024-02","configured_by":"admin"}`
	w2 := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data2))

	assert.Equal(t, http.StatusBadRequest, w2.Code)
	assert.Contains(t, w2.Body.String(), "utility rate for period already exists")
}

// ==================== UTILITY RATE GET ALL TESTS ====================

func TestE2E_UtilityRate_GetAll_Success(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, staffToken := registerUserForUtilityRateE2E("Staff User", "5555555555", "ur_staff4@test.com", "STAFF")
	_, tenantToken := registerUserForUtilityRateE2E("Tenant User", "6666666666", "ur_tenant2@test.com", "TENANT")

	// Create rates
	data1 := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-03","configured_by":"admin"}`
	MakeRequest("POST", "/utility-rates/", staffToken, strings.NewReader(data1))

	data2 := `{"water_rate":6.0,"electric_rate":8.0,"common_fee":3.0,"period":"2024-04","configured_by":"admin"}`
	MakeRequest("POST", "/utility-rates/", staffToken, strings.NewReader(data2))

	// Staff can get all
	w1 := MakeRequest("GET", "/utility-rates/", staffToken, nil)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Contains(t, w1.Body.String(), "Rates retrieved successfully")

	// Tenant can also get all
	w2 := MakeRequest("GET", "/utility-rates/", tenantToken, nil)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "Rates retrieved successfully")
}

func TestE2E_UtilityRate_GetAll_Unauthorized(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	w := MakeRequest("GET", "/utility-rates/", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_UtilityRate_GetAll_Empty(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "7777777777", "ur_staff5@test.com", "STAFF")

	w := MakeRequest("GET", "/utility-rates/", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Rates retrieved successfully")
}

// ==================== UTILITY RATE GET BY ID TESTS ====================

func TestE2E_UtilityRate_GetByID_Success(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, staffToken := registerUserForUtilityRateE2E("Staff User", "8888888888", "ur_staff6@test.com", "STAFF")

	// Create a rate
	createData := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-05","configured_by":"admin"}`
	createResp := MakeRequest("POST", "/utility-rates/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Extract ID from response (assuming response contains rate_id)
	responseBody := createResp.Body.String()
	// Find rate_id in response and extract it
	startIdx := strings.Index(responseBody, `"rate_id":"`)
	if startIdx == -1 {
		t.Fatal("Could not find rate_id in response")
	}
	rateIDStart := startIdx + 11
	endIdx := strings.Index(responseBody[rateIDStart:], `"`)
	rateID := responseBody[rateIDStart : rateIDStart+endIdx]

	// Get the rate by ID
	w := MakeRequest("GET", "/utility-rates/"+rateID, staffToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Rate retrieved successfully")
	assert.Contains(t, w.Body.String(), rateID)
}

func TestE2E_UtilityRate_GetByID_NotFound(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "9999999999", "ur_staff7@test.com", "STAFF")

	w := MakeRequest("GET", "/utility-rates/nonexistent-id", token, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Rate not found")
}

func TestE2E_UtilityRate_GetByID_Unauthorized(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	w := MakeRequest("GET", "/utility-rates/some-id", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ==================== UTILITY RATE UPDATE TESTS ====================

func TestE2E_UtilityRate_Update_Success(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1010101010", "ur_staff8@test.com", "STAFF")

	// Create a rate
	createData := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-06","configured_by":"admin"}`
	createResp := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Extract rate ID
	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"rate_id":"`)
	if startIdx == -1 {
		t.Fatal("Could not find rate_id in response")
	}
	rateIDStart := startIdx + 11
	endIdx := strings.Index(responseBody[rateIDStart:], `"`)
	rateID := responseBody[rateIDStart : rateIDStart+endIdx]

	// Update the rate
	updateData := `{"water_rate":6.0,"electric_rate":8.0,"common_fee":3.0}}`
	w := MakeRequest("PUT", "/utility-rates/"+rateID, token, strings.NewReader(updateData))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Utility rate updated successfully")
	assert.Contains(t, w.Body.String(), "6")
	assert.Contains(t, w.Body.String(), "8")
}

func TestE2E_UtilityRate_Update_NotFound(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1111101010", "ur_staff9@test.com", "STAFF")

	updateData := `{"water_rate":6.0,"electric_rate":8.0,"common_fee":3.0}}`
	w := MakeRequest("PUT", "/utility-rates/nonexistent-id", token, strings.NewReader(updateData))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to update rate")
}

func TestE2E_UtilityRate_Update_Unauthorized(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	updateData := `{"water_rate":6.0,"electric_rate":8.0,"common_fee":3.0}}`
	w := MakeRequest("PUT", "/utility-rates/some-id", "", strings.NewReader(updateData))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_UtilityRate_Update_Forbidden_Tenant(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, staffToken := registerUserForUtilityRateE2E("Staff User", "1212121212", "ur_staff10@test.com", "STAFF")
	_, tenantToken := registerUserForUtilityRateE2E("Tenant User", "1313131313", "ur_tenant3@test.com", "TENANT")

	// Create a rate as staff
	createData := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-07","configured_by":"admin"}`
	createResp := MakeRequest("POST", "/utility-rates/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Extract rate ID
	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"rate_id":"`)
	rateIDStart := startIdx + 11
	endIdx := strings.Index(responseBody[rateIDStart:], `"`)
	rateID := responseBody[rateIDStart : rateIDStart+endIdx]

	// Try to update as tenant
	updateData := `{"water_rate":6.0,"electric_rate":8.0,"common_fee":3.0}}`
	w := MakeRequest("PUT", "/utility-rates/"+rateID, tenantToken, strings.NewReader(updateData))

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==================== UTILITY RATE DELETE TESTS ====================

func TestE2E_UtilityRate_Delete_Success(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1414141414", "ur_staff11@test.com", "STAFF")

	// Create a rate
	createData := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-08","configured_by":"admin"}`
	createResp := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Extract rate ID
	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"rate_id":"`)
	rateIDStart := startIdx + 11
	endIdx := strings.Index(responseBody[rateIDStart:], `"`)
	rateID := responseBody[rateIDStart : rateIDStart+endIdx]

	// Delete the rate
	w := MakeRequest("DELETE", "/utility-rates/"+rateID, token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Utility rate deleted successfully")
}

func TestE2E_UtilityRate_Delete_NotFound(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1515151515", "ur_staff12@test.com", "STAFF")

	w := MakeRequest("DELETE", "/utility-rates/nonexistent-id", token, nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete rate")
}

func TestE2E_UtilityRate_Delete_Unauthorized(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	w := MakeRequest("DELETE", "/utility-rates/some-id", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestE2E_UtilityRate_Delete_Forbidden_Tenant(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, staffToken := registerUserForUtilityRateE2E("Staff User", "1616161616", "ur_staff13@test.com", "STAFF")
	_, tenantToken := registerUserForUtilityRateE2E("Tenant User", "1717171717", "ur_tenant4@test.com", "TENANT")

	// Create a rate as staff
	createData := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":2.5,"period":"2024-09","configured_by":"admin"}`
	createResp := MakeRequest("POST", "/utility-rates/", staffToken, strings.NewReader(createData))
	assert.Equal(t, http.StatusCreated, createResp.Code)

	// Extract rate ID
	responseBody := createResp.Body.String()
	startIdx := strings.Index(responseBody, `"rate_id":"`)
	rateIDStart := startIdx + 11
	endIdx := strings.Index(responseBody[rateIDStart:], `"`)
	rateID := responseBody[rateIDStart : rateIDStart+endIdx]

	// Try to delete as tenant
	w := MakeRequest("DELETE", "/utility-rates/"+rateID, tenantToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==================== UTILITY RATE INVALID DATA TESTS ====================

func TestE2E_UtilityRate_Create_InvalidJSON(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1818181818", "ur_staff14@test.com", "STAFF")

	data := `{invalid json}`
	w := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestE2E_UtilityRate_Create_NegativeElectricRate(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "1919191919", "ur_staff15@test.com", "STAFF")

	data := `{"water_rate":5.5,"electric_rate":-7.2,"common_fee":2.5,"period":"2024-10","configured_by":"admin"}`
	w := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid rates")
}

func TestE2E_UtilityRate_Create_NegativeCommonFee(t *testing.T) {
	setupUtilityRateTest()
	defer cleanupUtilityRateTest()

	_, token := registerUserForUtilityRateE2E("Staff User", "2020202020", "ur_staff16@test.com", "STAFF")

	data := `{"water_rate":5.5,"electric_rate":7.2,"common_fee":-2.5,"period":"2024-11","configured_by":"admin"}`
	w := MakeRequest("POST", "/utility-rates/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid rates")
}
