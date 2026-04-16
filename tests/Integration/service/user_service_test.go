package service

import (
	"testing"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanupTestUsers() {
	setup.CleanupUsers([]string{})
	setup.ResetTestDB()
}

// ==================== CREATE USER TESTS ====================

func TestUserService_CreateUser_EmptyName(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user := model.NewUser("", "1234567890", "emptyname@test.com", "password123", "TENANT")
	_, err := setup.UserService.CreateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_CreateUser_EmptyRole(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user := model.NewUser("Test User", "1234567890", "emptyrole@test.com", "password123", "")
	_, err := setup.UserService.CreateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_CreateUser_SpecialCharactersInName(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user := model.NewUser("John O'Brien", "1234567890", "special@test.com", "password123", "TENANT")
	createdUser, err := setup.UserService.CreateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "John O'Brien", createdUser.Name)
}

func TestUserService_CreateUser_ThaiCharacters(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user := model.NewUser("สมชาย ใจดี", "1234567890", "thai@test.com", "password123", "TENANT")
	createdUser, err := setup.UserService.CreateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "สมชาย ใจดี", createdUser.Name)
}

// ==================== UPDATE USER TESTS ====================

func TestUserService_UpdateUser_EmptyName(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Update Name", "1234567890", "updatename@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = ""
	user.Email = "updatename@test.com"

	_, err = setup.UserService.UpdateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_UpdateUser_EmptyEmail(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Update Email2", "1234567890", "updateemail2@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = "Updated Name"
	user.Email = ""

	_, err = setup.UserService.UpdateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_UpdateUser_KeepRoleWhenStaffUpdates(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Staff Keep Role", "1234567890", "staffkeep@test.com", "password123", "STAFF")
	require.NoError(t, err)

	originalRole := user.Role
	user.Name = "Updated Staff Name"

	updatedUser, err := setup.UserService.UpdateUser(user)

	require.NoError(t, err)
	assert.Equal(t, originalRole, updatedUser.Role)
	assert.Equal(t, "Updated Staff Name", updatedUser.Name)
}

func TestUserService_UpdateUser_UpdatePhone(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Update Phone", "1234567890", "updatephone@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Phone = "9998887777"
	updatedUser, err := setup.UserService.UpdateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "9998887777", updatedUser.Phone)
}

func TestUserService_UpdateUser_MultipleUpdates(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Multi Update", "1234567890", "multiupdate@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = "First Update"
	user.Phone = "1111111111"
	_, err = setup.UserService.UpdateUser(user)
	require.NoError(t, err)

	user.Name = "Second Update"
	user.Phone = "2222222222"
	updatedUser, err := setup.UserService.UpdateUser(user)
	require.NoError(t, err)

	assert.Equal(t, "Second Update", updatedUser.Name)
	assert.Equal(t, "2222222222", updatedUser.Phone)
}

// ==================== GET USER TESTS ====================

func TestUserService_GetUserByID_StaffCanViewAnyUser(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	anyUser, err := setup.AuthService.Register("Any User", "1234567890", "anyuser@test.com", "userpass", "TENANT")
	require.NoError(t, err)

	user, err := setup.UserService.GetUserByID(anyUser.ID)

	require.NoError(t, err)
	assert.Equal(t, anyUser.ID, user.ID)
}

func TestUserService_GetUserByID_VerifyUserFields(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Verify Fields", "0987654321", "verifyfields@test.com", "password123", "STAFF")
	require.NoError(t, err)

	retrievedUser, err := setup.UserService.GetUserByID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, "Verify Fields", retrievedUser.Name)
	assert.Equal(t, "0987654321", retrievedUser.Phone)
	assert.Equal(t, "verifyfields@test.com", retrievedUser.Email)
	assert.Equal(t, "STAFF", retrievedUser.Role)
}

// ==================== GET USERS BY ROLE TESTS ====================

func TestUserService_GetUsersByRole_OnlyReturnsMatchingRole(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	_, err := setup.AuthService.Register("Role Check 1", "1234567890", "rolecheck1@test.com", "pass1", "STAFF")
	require.NoError(t, err)

	_, err = setup.AuthService.Register("Role Check 2", "1234567890", "rolecheck2@test.com", "pass2", "TENANT")
	require.NoError(t, err)

	_, err = setup.AuthService.Register("Role Check 3", "1234567890", "rolecheck3@test.com", "pass3", "TENANT")
	require.NoError(t, err)

	tenants, err := setup.UserService.GetUsersByRole("TENANT")

	require.NoError(t, err)
	for _, tenant := range tenants {
		assert.Equal(t, "TENANT", tenant.Role)
	}
}

func TestUserService_GetUsersByRole_CaseSensitivity(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	_, err := setup.AuthService.Register("Case Sen", "1234567890", "casesen@test.com", "pass1", "TENANT")
	require.NoError(t, err)

	uppercaseResult, err := setup.UserService.GetUsersByRole("TENANT")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(uppercaseResult), 1)
}

// ==================== DELETE USER TESTS ====================

func TestUserService_DeleteUser_AfterUpdate(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Delete After Update", "1234567890", "deleteafterupdate@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = "Updated Before Delete"
	_, err = setup.UserService.UpdateUser(user)
	require.NoError(t, err)

	err = setup.UserService.DeleteUser(user.ID)
	require.NoError(t, err)

	_, err = setup.UserService.GetUserByID(user.ID)
	assert.Error(t, err)
}

func TestUserService_DeleteUser_AlreadyDeleted(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user, err := setup.AuthService.Register("Double Delete", "1234567890", "double@test.com", "password123", "TENANT")
	require.NoError(t, err)

	err = setup.UserService.DeleteUser(user.ID)
	require.NoError(t, err)

	err = setup.UserService.DeleteUser(user.ID)
	assert.Error(t, err)
}

// ==================== AUTH SERVICE TESTS ====================

func TestAuthService_Register_PhoneNumberFormats(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	user1, err := setup.AuthService.Register("Phone 1", "0861234567", "phone1@test.com", "pass1", "TENANT")
	require.NoError(t, err)
	assert.Equal(t, "0861234567", user1.Phone)

	user2, err := setup.AuthService.Register("Phone 2", "+66812345678", "phone2@test.com", "pass2", "TENANT")
	require.NoError(t, err)
	assert.Equal(t, "+66812345678", user2.Phone)
}

func TestAuthService_Register_VeryLongEmail(t *testing.T) {
	setup.ResetTestDB()
	defer cleanupTestUsers()

	longEmail := "verylongemailwithlotsofcharacters1234567890@test.com"
	user, err := setup.AuthService.Register("Long Email", "1234567890", longEmail, "password123", "TENANT")

	require.NoError(t, err)
	assert.Equal(t, longEmail, user.Email)
}
