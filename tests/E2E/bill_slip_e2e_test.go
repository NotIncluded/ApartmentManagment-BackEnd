package e2e

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"io"

	// "github.com/PunMung-66/ApartmentSys/internal/auth"
	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	// "strings"
)
// var secret = []byte(setup.JWTSecret)


// MakeRequestWithBody sends a multipart/form-data request for testing
func MakeRequestWithBody(method, path, token string, body *bytes.Buffer, contentType string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if token != "" {
		 req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	// testRouter is defined in e2e_test.go, so we need to initialize it here if not already
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

func createTestBillAndRoom() (string, string) {
	room := setup.CreateTestRoom("S101", 1, "Available")
	bill := setup.CreateTestBill(room.ID)
	return bill.ID, room.ID
}

func TestE2E_BillSlip_Upload_Success(t *testing.T) {
	setupBillSlipTest()
	defer cleanupBillSlipTest()

	_, staffToken := registerUserForSlipE2E("Staff User", "1111111111", "slip_staff@test.com", "STAFF")
	billID, roomID := createTestBillAndRoom()

	filePath := "testdata/test-slip.jpg" // Place a test image at this path
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_ = w.WriteField("bill_id", billID)
	_ = w.WriteField("room_id", roomID)
	fw, err := w.CreateFormFile("slip", filepath.Base(filePath))
	assert.NoError(t, err)
	_, err = file.Seek(0, 0)
	assert.NoError(t, err)
	_, err = io.Copy(fw, file)
	assert.NoError(t, err)
	w.Close()

	resp := MakeRequestWithBody("POST", "/billslips/upload", staffToken, &b, w.FormDataContentType())
	// GenerateTestToken is imported from e2e_test.go
	// MakeRequestWithBody is imported from e2e_test.go
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "slip_url")
}
