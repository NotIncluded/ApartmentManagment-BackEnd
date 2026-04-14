package auth_user_service_integration

import (
	"testing"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser_EmptyName(t *testing.T) {
	user := model.NewUser("", "1234567890", "emptyname@test.com", "password123", "TENANT")

	_, err := userService.CreateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_CreateUser_EmptyRole(t *testing.T) {
	user := model.NewUser("Test User", "1234567890", "emptyrole@test.com", "password123", "")

	_, err := userService.CreateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_UpdateUser_EmptyName(t *testing.T) {
	defer cleanupTestUsers([]string{"updatename@test.com"})

	user, err := authService.Register("Update Name", "1234567890", "updatename@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = ""
	user.Email = "updatename@test.com"

	_, err = userService.UpdateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_UpdateUser_EmptyEmail(t *testing.T) {
	defer cleanupTestUsers([]string{"updateemail2@test.com"})

	user, err := authService.Register("Update Email2", "1234567890", "updateemail2@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = "Updated Name"
	user.Email = ""

	_, err = userService.UpdateUser(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete request body")
}

func TestUserService_UpdateUser_KeepRoleWhenStaffUpdates(t *testing.T) {
	defer cleanupTestUsers([]string{"staffkeep@test.com"})

	user, err := authService.Register("Staff Keep Role", "1234567890", "staffkeep@test.com", "password123", "STAFF")
	require.NoError(t, err)

	originalRole := user.Role
	user.Name = "Updated Staff Name"

	updatedUser, err := userService.UpdateUser(user)

	require.NoError(t, err)
	assert.Equal(t, originalRole, updatedUser.Role)
	assert.Equal(t, "Updated Staff Name", updatedUser.Name)
}

func TestUserService_GetUserByID_StaffCanViewAnyUser(t *testing.T) {
	defer cleanupTestUsers([]string{"staffview@test.com", "anyuser@test.com"})

	anyUser, err := authService.Register("Any User", "1234567890", "anyuser@test.com", "userpass", "TENANT")
	require.NoError(t, err)

	user, err := userService.GetUserByID(anyUser.ID)

	require.NoError(t, err)
	assert.Equal(t, anyUser.ID, user.ID)
}

func TestUserService_GetUserByID_VerifyUserFields(t *testing.T) {
	defer cleanupTestUsers([]string{"verifyfields@test.com"})

	user, err := authService.Register("Verify Fields", "0987654321", "verifyfields@test.com", "password123", "STAFF")
	require.NoError(t, err)

	retrievedUser, err := userService.GetUserByID(user.ID)

	require.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, "Verify Fields", retrievedUser.Name)
	assert.Equal(t, "0987654321", retrievedUser.Phone)
	assert.Equal(t, "verifyfields@test.com", retrievedUser.Email)
	assert.Equal(t, "STAFF", retrievedUser.Role)
}

func TestUserService_GetUsersByRole_OnlyReturnsMatchingRole(t *testing.T) {
	defer cleanupTestUsers([]string{"rolecheck1@test.com", "rolecheck2@test.com", "rolecheck3@test.com"})

	_, err := authService.Register("Role Check 1", "1234567890", "rolecheck1@test.com", "pass1", "STAFF")
	require.NoError(t, err)

	_, err = authService.Register("Role Check 2", "1234567890", "rolecheck2@test.com", "pass2", "TENANT")
	require.NoError(t, err)

	_, err = authService.Register("Role Check 3", "1234567890", "rolecheck3@test.com", "pass3", "TENANT")
	require.NoError(t, err)

	tenants, err := userService.GetUsersByRole("TENANT")

	require.NoError(t, err)
	for _, tenant := range tenants {
		assert.Equal(t, "TENANT", tenant.Role)
	}
}

func TestUserService_DeleteUser_AfterUpdate(t *testing.T) {
	defer cleanupTestUsers([]string{"deleteafterupdate@test.com"})

	user, err := authService.Register("Delete After Update", "1234567890", "deleteafterupdate@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = "Updated Before Delete"
	_, err = userService.UpdateUser(user)
	require.NoError(t, err)

	err = userService.DeleteUser(user.ID)
	require.NoError(t, err)

	_, err = userService.GetUserByID(user.ID)
	assert.Error(t, err)
}

func TestUserService_CreateUser_SpecialCharactersInName(t *testing.T) {
	defer cleanupTestUsers([]string{"special@test.com"})

	user := model.NewUser("John O'Brien", "1234567890", "special@test.com", "password123", "TENANT")

	createdUser, err := userService.CreateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "John O'Brien", createdUser.Name)
}

func TestUserService_CreateUser_ThaiCharacters(t *testing.T) {
	defer cleanupTestUsers([]string{"thai@test.com"})

	user := model.NewUser("สมชาย ใจดี", "1234567890", "thai@test.com", "password123", "TENANT")

	createdUser, err := userService.CreateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "สมชาย ใจดี", createdUser.Name)
}

func TestAuthService_Register_PhoneNumberFormats(t *testing.T) {
	defer cleanupTestUsers([]string{"phone1@test.com", "phone2@test.com"})

	user1, err := authService.Register("Phone 1", "0861234567", "phone1@test.com", "pass1", "TENANT")
	require.NoError(t, err)
	assert.Equal(t, "0861234567", user1.Phone)

	user2, err := authService.Register("Phone 2", "+66812345678", "phone2@test.com", "pass2", "TENANT")
	require.NoError(t, err)
	assert.Equal(t, "+66812345678", user2.Phone)
}

func TestUserService_UpdateUser_UpdatePhone(t *testing.T) {
	defer cleanupTestUsers([]string{"updatephone@test.com"})

	user, err := authService.Register("Update Phone", "1234567890", "updatephone@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Phone = "9998887777"
	updatedUser, err := userService.UpdateUser(user)

	require.NoError(t, err)
	assert.Equal(t, "9998887777", updatedUser.Phone)
}

func TestUserService_GetUsersByRole_CaseSensitivity(t *testing.T) {
	defer cleanupTestUsers([]string{"casesen@test.com"})

	_, err := authService.Register("Case Sen", "1234567890", "casesen@test.com", "pass1", "TENANT")
	require.NoError(t, err)

	uppercaseResult, err := userService.GetUsersByRole("TENANT")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(uppercaseResult), 1)
}

func TestUserService_DeleteUser_AlreadyDeleted(t *testing.T) {
	defer cleanupTestUsers([]string{"double@test.com"})

	user, err := authService.Register("Double Delete", "1234567890", "double@test.com", "password123", "TENANT")
	require.NoError(t, err)

	err = userService.DeleteUser(user.ID)
	require.NoError(t, err)

	err = userService.DeleteUser(user.ID)
	assert.Error(t, err)
}

func TestAuthService_Register_VeryLongEmail(t *testing.T) {
	defer cleanupTestUsers([]string{"verylongemailwithlotsofcharacters1234567890@test.com"})

	longEmail := "verylongemailwithlotsofcharacters1234567890@test.com"
	user, err := authService.Register("Long Email", "1234567890", longEmail, "password123", "TENANT")

	require.NoError(t, err)
	assert.Equal(t, longEmail, user.Email)
}

func TestUserService_UpdateUser_MultipleUpdates(t *testing.T) {
	defer cleanupTestUsers([]string{"multiupdate@test.com"})

	user, err := authService.Register("Multi Update", "1234567890", "multiupdate@test.com", "password123", "TENANT")
	require.NoError(t, err)

	user.Name = "First Update"
	user.Phone = "1111111111"
	_, err = userService.UpdateUser(user)
	require.NoError(t, err)

	user.Name = "Second Update"
	user.Phone = "2222222222"
	updatedUser, err := userService.UpdateUser(user)
	require.NoError(t, err)

	assert.Equal(t, "Second Update", updatedUser.Name)
	assert.Equal(t, "2222222222", updatedUser.Phone)
}
