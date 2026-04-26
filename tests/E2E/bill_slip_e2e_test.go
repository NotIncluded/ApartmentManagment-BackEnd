package e2e

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
)

// MakeRequestWithBody sends a multipart/form-data request for testing
func MakeRequestWithBody(method, path, token string, body *bytes.Buffer, contentType string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	
	if testRouter == nil {
		panic("testRouter is not initialized. Please run this test with other E2E tests or refactor helpers.")
	}
	testRouter.ServeHTTP(w, req)
	return w
}

func setupBillSlipTest() {
	setup.ResetTestDB()
}

func cleanupBillSlipTest() {
	setup.ResetTestDB()
}

func registerUserForSlipE2E(name, phone, email, role string) (*setup.User, string) {
	user, _ := setup.AuthService.Register(name, phone, email, "password123", role)
	token := GenerateTestToken(user.ID, role)
	return user, token
}

func createTestBillAndRoom(userID string) (string, string) {
	room := setup.CreateTestRoom("S101", 1, "Available")
	bill := setup.CreateTestBill(userID, room.ID)
	return bill.ID, room.ID
}

func TestE2E_BillSlip_Upload_Success(t *testing.T) {
	setupBillSlipTest()
	defer cleanupBillSlipTest()

	staffUser, staffToken := registerUserForSlipE2E("Staff User", "1111111111", "slip_staff@test.com", "STAFF")
	billID, roomID := createTestBillAndRoom(staffUser.ID)

	// 1. Create a memory buffer to hold our fake form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 2. Write the standard text fields
	_ = w.WriteField("bill_id", billID)
	_ = w.WriteField("room_id", roomID)

	// 3. Create a fake file completely in memory instead of reading from disk
	fw, err := w.CreateFormFile("slip", "mock-slip.jpg")
	assert.NoError(t, err)
	
	// Write a few bytes of dummy data to simulate an image file
	_, err = fw.Write([]byte("fake-image-bytes-for-testing-purposes"))
	assert.NoError(t, err)
	
	// 4. Close the writer BEFORE making the request to finalize the multipart boundaries
	w.Close()

	// 5. Fire the request
	resp := MakeRequestWithBody("POST", "/billslips/upload", staffToken, &b, w.FormDataContentType())

	// 6. Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "slip_url")
}