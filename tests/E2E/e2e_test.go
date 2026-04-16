package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PunMung-66/ApartmentSys/controller"
	"github.com/PunMung-66/ApartmentSys/internal/auth"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	testRouter *gin.Engine
	secret     = []byte(setup.JWTSecret)
)

func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	userRepo := repository.NewUserRepository(setup.TestDB)
	roomRepo := repository.NewRoomRepository(setup.TestDB)
	contractRepo := repository.NewContractRepository(setup.TestDB)

	userService := service.NewUserService(userRepo)
	roomService := service.NewRoomService(roomRepo, contractRepo)
	roomService.SetUserRepository(userRepo)
	authService := service.NewAuthService(userRepo)
	contractService := service.NewContractService(contractRepo, roomRepo)
	contractService.SetUserRepository(userRepo)

	userController := controller.NewUserController(userService)
	roomController := controller.NewRoomController(roomService)
	contractController := controller.NewContractController(contractService)
	authController := controller.NewAuthController(authService, secret)

	userRoute := r.Group("/users")
	{
		userRoute.POST("/", auth.Protect(secret, "STAFF"), userController.CreateUser)
		userRoute.GET("/", auth.Protect(secret, "STAFF"), userController.GetUsersByRole)
		userRoute.GET("/:id", auth.Protect(secret, "STAFF", "TENANT"), userController.GetUserByID)
		userRoute.PUT("/:id", auth.Protect(secret, "STAFF", "TENANT"), userController.UpdateUser)
		userRoute.DELETE("/:id", auth.Protect(secret, "STAFF"), userController.DeleteUser)
	}

	roomRoute := r.Group("/rooms")
	{
		roomRoute.POST("/", auth.Protect(secret, "STAFF"), roomController.CreateRoom)
		roomRoute.GET("/", auth.Protect(secret, "STAFF"), roomController.GetListRoom)
		roomRoute.GET("/:id", auth.Protect(secret, "STAFF"), roomController.GetRoomByID)
		roomRoute.PUT("/:id", auth.Protect(secret, "STAFF"), roomController.UpdateRoom)
		roomRoute.DELETE("/:id", auth.Protect(secret, "STAFF"), roomController.DeleteRoom)
		roomRoute.GET("/:id/contract", auth.Protect(secret, "STAFF"), roomController.GetRoomActiveContract)
		roomRoute.GET("/:id/contracts", auth.Protect(secret, "STAFF"), roomController.GetRoomContractHistory)
		roomRoute.GET("/:id/tenant", auth.Protect(secret, "STAFF"), roomController.GetRoomTenant)
		roomRoute.POST("/:id/assign", auth.Protect(secret, "STAFF"), roomController.AssignRoom)
	}

	contractRoute := r.Group("/contracts")
	{
		contractRoute.POST("/", auth.Protect(secret, "STAFF"), contractController.CreateContract)
		contractRoute.GET("/", auth.Protect(secret, "STAFF"), contractController.GetContracts)
		contractRoute.GET("/:id", auth.Protect(secret, "STAFF"), contractController.GetContractByID)
		contractRoute.PUT("/:id", auth.Protect(secret, "STAFF"), contractController.UpdateContract)
		contractRoute.DELETE("/:id", auth.Protect(secret, "STAFF"), contractController.DeleteContract)
		contractRoute.GET("/user/:userID", auth.Protect(secret, "STAFF"), contractController.GetContractsByUserID)
		contractRoute.GET("/room/:roomID", auth.Protect(secret, "STAFF"), contractController.GetContractsByRoomID)
	}

	meRoute := r.Group("/me")
	{
		meRoute.GET("/room", auth.Protect(secret, "TENANT"), roomController.GetMyRoom)
	}

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", authController.LoginHandler)
		authRoute.POST("/register", authController.RegisterHandler)
	}

	return r
}

func GenerateTestToken(userID, role string) string {
	token, _ := auth.GenerateToken(secret, userID, role)
	return token
}

func MakeRequest(method, path string, token string, body *strings.Reader) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, path, body)
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func init() {
	testRouter = SetupTestRouter()
}

func setupAuthTest() {
	setup.ResetTestDB()
}

func cleanupAuthTest() {
	setup.ResetTestDB()
}

// ==================== AUTH REGISTER TESTS ====================

func TestE2E_Auth_Register_Success(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"name":"Register Test","phone":"1111111111","email":"register@test.com","password":"password123","role":"TENANT"}`
	w := MakeRequest("POST", "/auth/register", "", strings.NewReader(data))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
	assert.Contains(t, w.Body.String(), "register@test.com")
}

func TestE2E_Auth_Register_DuplicateEmail(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	setup.AuthService.Register("First User", "1111111111", "duplicate@test.com", "password123", "TENANT")

	data := `{"name":"Second User","phone":"2222222222","email":"duplicate@test.com","password":"password456","role":"TENANT"}`
	w := MakeRequest("POST", "/auth/register", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Registration failed")
}

func TestE2E_Auth_Register_MissingFields(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"name":"","phone":"1111111111","email":"","password":"","role":"TENANT"}`
	w := MakeRequest("POST", "/auth/register", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestE2E_Auth_Register_InvalidEmail(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"name":"Invalid Email","phone":"1111111111","email":"not-an-email","password":"password123","role":"TENANT"}`
	w := MakeRequest("POST", "/auth/register", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestE2E_Auth_Register_CannotSetStaffRole(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"name":"Staff Hack","phone":"1111111111","email":"staffhack@test.com","password":"password123","role":"STAFF"}`
	w := MakeRequest("POST", "/auth/register", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Cannot set STAFF role")
}

// ==================== AUTH LOGIN TESTS ====================

func TestE2E_Auth_Login_Success(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	setup.AuthService.Register("Login User", "1111111111", "login@test.com", "password123", "TENANT")

	data := `{"email":"login@test.com","password":"password123"}`
	w := MakeRequest("POST", "/auth/login", "", strings.NewReader(data))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Access token generated successfully")
	assert.Contains(t, w.Body.String(), "access_token")
}

func TestE2E_Auth_Login_InvalidPassword(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	setup.AuthService.Register("Wrong Pass User", "1111111111", "wrongpass@test.com", "password123", "TENANT")

	data := `{"email":"wrongpass@test.com","password":"wrongpassword"}`
	w := MakeRequest("POST", "/auth/login", "", strings.NewReader(data))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestE2E_Auth_Login_InvalidEmail(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"email":"nonexistent@test.com","password":"password123"}`
	w := MakeRequest("POST", "/auth/login", "", strings.NewReader(data))

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}

func TestE2E_Auth_Login_MissingEmail(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"email":"","password":"password123"}`
	w := MakeRequest("POST", "/auth/login", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestE2E_Auth_Login_MissingPassword(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{"email":"login@test.com","password":""}`
	w := MakeRequest("POST", "/auth/login", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestE2E_Auth_Login_EmptyBody(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	data := `{}`
	w := MakeRequest("POST", "/auth/login", "", strings.NewReader(data))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ==================== USER CRUD TESTS ====================

func TestE2E_User_Create_Success(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	_, token := registerUserForE2E("Staff User", "1111111111", "e2estaff@test.com", "STAFF")

	data := `{"name":"New User","phone":"2222222222","email":"newuser@test.com","password":"password123","role":"TENANT"}`
	w := MakeRequest("POST", "/users/", token, strings.NewReader(data))

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Add user successfully")
}

func TestE2E_User_GetAll_Tenants(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	_, token := registerUserForE2E("Staff User", "1111111111", "e2estaff@test.com", "STAFF")
	setup.AuthService.Register("Tenant 1", "3333333333", "tenant1@test.com", "pass1", "TENANT")
	setup.AuthService.Register("Tenant 2", "4444444444", "tenant2@test.com", "pass2", "TENANT")

	w := MakeRequest("GET", "/users/?role=TENANT", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Users retrieved successfully")
}

func TestE2E_User_GetByID_Success(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	_, token := registerUserForE2E("Staff User", "1111111111", "e2estaff@test.com", "STAFF")
	user, _ := setup.AuthService.Register("Get User", "5555555555", "getuser@test.com", "password123", "TENANT")

	w := MakeRequest("GET", "/users/"+user.ID, token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User retrieved successfully")
}

func TestE2E_User_Update_Success(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	_, token := registerUserForE2E("Staff User", "1111111111", "e2estaff@test.com", "STAFF")
	user, _ := setup.AuthService.Register("Update User", "6666666666", "updateuser@test.com", "password123", "TENANT")

	data := `{"name":"Updated Name","phone":"7777777777","email":"updateuser@test.com","password":"password123","role":"TENANT"}`
	w := MakeRequest("PUT", "/users/"+user.ID, token, strings.NewReader(data))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User updated successfully")
}

func TestE2E_User_Delete_Success(t *testing.T) {
	setupAuthTest()
	defer cleanupAuthTest()

	_, token := registerUserForE2E("Staff User", "1111111111", "e2estaff@test.com", "STAFF")
	user, _ := setup.AuthService.Register("Delete User", "8888888888", "deleteuser@test.com", "password123", "TENANT")

	w := MakeRequest("DELETE", "/users/"+user.ID, token, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User deleted successfully")
}
