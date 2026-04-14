package integration

import (
	"os"
	"testing"

	"github.com/PunMung-66/ApartmentSys/config"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	testDBName     = "apartment_test"
	jwtSecret      = "test_jwt_secret_key"
	testDBHost     = "3.80.58.141"
	testDBUser     = "postgres"
	testDBPassword = "mysecretpassword"
	testDBPort     = "5432"
)

var (
	testDB      *gorm.DB
	userRepo    *repository.UserRepository
	authService *service.AuthService
	userService *service.UserService
)

func TestMain(m *testing.M) {
	setupTestDB()
	resetTestDB()
	runTests := m.Run()
	// teardownTestDB()
	os.Exit(runTests)
}

func setupTestDB() {
	os.Setenv("DB_HOST", testDBHost)
	os.Setenv("DB_USER", testDBUser)
	os.Setenv("DB_PASSWORD", testDBPassword)
	os.Setenv("DB_PORT", testDBPort)

	db, err := config.ConnectTestDatabase(testDBName)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	db.AutoMigrate(&model.User{})

	testDB = db
	userRepo = repository.NewUserRepository(testDB)
	authService = service.NewAuthService(userRepo)
	userService = service.NewUserService(userRepo)
}

func resetTestDB() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE users CASCADE")
	}
}

func teardownTestDB() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE users CASCADE")
		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func cleanupTestUsers(emails []string) {
	for _, email := range emails {
		testDB.Unscoped().Where("email = ?", email).Delete(&model.User{})
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	defer cleanupTestUsers([]string{"newuser@test.com"})

	user, err := authService.Register("New User", "1234567890", "newuser@test.com", "password123", "TENANT")

	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "newuser@test.com", user.Email)
	assert.Equal(t, "TENANT", user.Role)
	assert.Equal(t, "New User", user.Name)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	defer cleanupTestUsers([]string{"duplicate@test.com"})

	_, err := authService.Register("First User", "1234567890", "duplicate@test.com", "password123", "TENANT")
	require.NoError(t, err)

	_, err = authService.Register("Second User", "9876543210", "duplicate@test.com", "password456", "TENANT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestAuthService_Login_Success(t *testing.T) {
	defer cleanupTestUsers([]string{"logintest@test.com"})

	_, err := authService.Register("Login Test", "1234567890", "logintest@test.com", "password123", "TENANT")
	require.NoError(t, err)

	token, err := authService.Login(service.LoginRequest{
		Email:    "logintest@test.com",
		Password: "password123",
	}, []byte(jwtSecret))

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	_, err := authService.Login(service.LoginRequest{
		Email:    "nonexistent@test.com",
		Password: "password123",
	}, []byte(jwtSecret))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email or password")
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	defer cleanupTestUsers([]string{"wrongpass@test.com"})

	_, err := authService.Register("Wrong Pass", "1234567890", "wrongpass@test.com", "correctpassword", "TENANT")
	require.NoError(t, err)

	_, err = authService.Login(service.LoginRequest{
		Email:    "wrongpass@test.com",
		Password: "wrongpassword",
	}, []byte(jwtSecret))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email or password")
}

func TestAuthService_Register_StaffRole(t *testing.T) {
	// defer cleanupTestUsers([]string{"staffuser@test.com"})

	user, err := authService.Register("Staff User", "1234567890", "staffuser@test.com", "staffpass", "STAFF")

	require.NoError(t, err)
	assert.Equal(t, "STAFF", user.Role)
}

func TestAuthService_Register_TenantRole(t *testing.T) {
	defer cleanupTestUsers([]string{"tenantuser@test.com"})

	user, err := authService.Register("Tenant User", "1234567890", "tenantuser@test.com", "tenantpass", "TENANT")

	require.NoError(t, err)
	assert.Equal(t, "TENANT", user.Role)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	defer cleanupTestUsers([]string{"servicecreate@test.com"})

	user := model.NewUser("Service Create", "1234567890", "servicecreate@test.com", "password123", "TENANT")

	createdUser, err := userService.CreateUser(user)

	require.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, "servicecreate@test.com", createdUser.Email)
}

func TestUserService_CreateUser_IncompleteRequest(t *testing.T) {
	user := model.NewUser("", "1234567890", "incomplete@test.com", "password123", "TENANT")

	_, err := userService.CreateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_CreateUser_StaffUser(t *testing.T) {
	defer cleanupTestUsers([]string{"staffbystaff@test.com"})

	staffUser := model.NewUser("Staff", "1234567890", "staffbystaff@test.com", "staffpass", "STAFF")

	createdUser, err := userService.CreateUser(staffUser)

	require.NoError(t, err)
	assert.Equal(t, "STAFF", createdUser.Role)
}

func TC10_TestUserService_DeleteUser_Success(t *testing.T) {
	defer cleanupTestUsers([]string{"deleteuser@test.com"})

	user, err := authService.Register("Delete User", "1234567890", "deleteuser@test.com", "password123", "TENANT")
	require.NoError(t, err)

	err = userService.DeleteUser(user.ID)
	require.NoError(t, err)

	_, err = userRepo.FindUserByEmail(&user.Email)
	assert.Error(t, err)
}

func TC11_TestUserService_DeleteUser_NotFound(t *testing.T) {
	err := userService.DeleteUser("non-existent-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TC12_TestAuthService_Login_StaffUser(t *testing.T) {
	defer cleanupTestUsers([]string{"stafflogin@test.com"})

	_, err := authService.Register("Staff Login", "1234567890", "stafflogin@test.com", "staff123", "STAFF")
	require.NoError(t, err)

	token, err := authService.Login(service.LoginRequest{
		Email:    "stafflogin@test.com",
		Password: "staff123",
	}, []byte(jwtSecret))

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TC13_TestPermissionRestrictions_StaffVsTenant(t *testing.T) {
	defer cleanupTestUsers([]string{"staffperm@test.com", "tenantperm@test.com"})

	staffUser, err := authService.Register("Staff Perm", "1234567890", "staffperm@test.com", "staffpass", "STAFF")
	require.NoError(t, err)
	assert.Equal(t, "STAFF", staffUser.Role)

	tenantUser, err := authService.Register("Tenant Perm", "1234567890", "tenantperm@test.com", "tenantpass", "TENANT")
	require.NoError(t, err)
	assert.Equal(t, "TENANT", tenantUser.Role)

	staffToken, err := authService.Login(service.LoginRequest{
		Email:    "staffperm@test.com",
		Password: "staffpass",
	}, []byte(jwtSecret))
	require.NoError(t, err)
	assert.NotEmpty(t, staffToken)

	tenantToken, err := authService.Login(service.LoginRequest{
		Email:    "tenantperm@test.com",
		Password: "tenantpass",
	}, []byte(jwtSecret))
	require.NoError(t, err)
	assert.NotEmpty(t, tenantToken)

	assert.NotEqual(t, staffToken, tenantToken)
}

func TC14_TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	defer cleanupTestUsers([]string{"servicedup@test.com"})

	user1 := model.NewUser("First", "1234567890", "servicedup@test.com", "pass1", "TENANT")
	_, err := userService.CreateUser(user1)
	require.NoError(t, err)

	user2 := model.NewUser("Second", "9876543210", "servicedup@test.com", "pass2", "TENANT")
	_, err = userService.CreateUser(user2)

	assert.Error(t, err)
}

func TC15_TestUserService_MultipleTenants(t *testing.T) {
	defer cleanupTestUsers([]string{"tenant1@test.com", "tenant2@test.com", "tenant3@test.com"})

	tenants := []struct {
		name  string
		email string
	}{
		{"Tenant One", "tenant1@test.com"},
		{"Tenant Two", "tenant2@test.com"},
		{"Tenant Three", "tenant3@test.com"},
	}

	for _, tenant := range tenants {
		user := model.NewUser(tenant.name, "1234567890", tenant.email, "password123", "TENANT")
		createdUser, err := userService.CreateUser(user)
		require.NoError(t, err, "Failed to create tenant: %s", tenant.email)
		assert.NotEmpty(t, createdUser.ID)
		assert.Equal(t, tenant.email, createdUser.Email)
	}
}
