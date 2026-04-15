package repository

import (
	"testing"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userRepositoryTest *repository.UserRepository
)

func init() {
	userRepositoryTest = repository.NewUserRepository(setup.TestDB)
}

// TestUserRepository_CreateUser_Success tests successful user creation
func TestUserRepository_CreateUser_Success(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user := model.NewUser("John Doe", "1234567890", "johndoe@test.com", "password123", "TENANT")

	createdUser, err := userRepositoryTest.CreateUser(user)

	require.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.Equal(t, "johndoe@test.com", createdUser.Email)
	assert.Equal(t, "1234567890", createdUser.Phone)
	assert.Equal(t, "TENANT", createdUser.Role)

	defer setup.CleanupUsers([]string{"johndoe@test.com"})
}

// TestUserRepository_CreateUser_MultipleUsers tests creating multiple users
func TestUserRepository_CreateUser_MultipleUsers(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user1 := model.NewUser("User One", "1111111111", "userone@test.com", "pass123", "TENANT")
	user2 := model.NewUser("User Two", "2222222222", "usertwo@test.com", "pass456", "STAFF")

	createdUser1, err := userRepositoryTest.CreateUser(user1)
	require.NoError(t, err)

	createdUser2, err := userRepositoryTest.CreateUser(user2)
	require.NoError(t, err)

	assert.NotEmpty(t, createdUser1.ID)
	assert.NotEmpty(t, createdUser2.ID)
	assert.NotEqual(t, createdUser1.ID, createdUser2.ID)
	assert.Equal(t, "User One", createdUser1.Name)
	assert.Equal(t, "User Two", createdUser2.Name)

	defer setup.CleanupUsers([]string{"userone@test.com", "usertwo@test.com"})
}

// TestUserRepository_FindUserByID_Success tests finding user by ID
func TestUserRepository_FindUserByID_Success(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user := model.NewUser("Find By ID", "3333333333", "findbyid@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	foundUser, err := userRepositoryTest.FindUserByID(createdUser.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, "Find By ID", foundUser.Name)
	assert.Equal(t, "findbyid@test.com", foundUser.Email)

	defer setup.CleanupUsers([]string{"findbyid@test.com"})
}

// TestUserRepository_FindUserByID_NotFound tests when user ID doesn't exist
func TestUserRepository_FindUserByID_NotFound(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	foundUser, err := userRepositoryTest.FindUserByID("nonexistent-user-id")

	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

// TestUserRepository_FindUserByEmail_Success tests finding user by email
func TestUserRepository_FindUserByEmail_Success(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user := model.NewUser("Email User", "4444444444", "emailuser@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	email := "emailuser@test.com"
	foundUser, err := userRepositoryTest.FindUserByEmail(&email)

	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, email, foundUser.Email)
	assert.Equal(t, "Email User", foundUser.Name)

	defer setup.CleanupUsers([]string{"emailuser@test.com"})
}

// TestUserRepository_FindUserByEmail_NotFound tests when email doesn't exist
func TestUserRepository_FindUserByEmail_NotFound(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	email := "nonexistent@test.com"
	foundUser, err := userRepositoryTest.FindUserByEmail(&email)

	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

// TestUserRepository_FindUsersByRole_Success tests finding all users by role
func TestUserRepository_FindUsersByRole_Success(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user1 := model.NewUser("Tenant One", "6666666666", "tenant1@test.com", "pass123", "TENANT")
	user2 := model.NewUser("Tenant Two", "7777777777", "tenant2@test.com", "pass456", "TENANT")
	user3 := model.NewUser("Staff One", "8888888888", "staff1@test.com", "pass789", "STAFF")

	_, err := userRepositoryTest.CreateUser(user1)
	require.NoError(t, err)

	_, err = userRepositoryTest.CreateUser(user2)
	require.NoError(t, err)

	_, err = userRepositoryTest.CreateUser(user3)
	require.NoError(t, err)

	tenantUsers, err := userRepositoryTest.FindUsersByRole("TENANT")

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(tenantUsers), 2)
	for _, u := range tenantUsers {
		assert.Equal(t, "TENANT", u.Role)
	}

	defer setup.CleanupUsers([]string{"tenant1@test.com", "tenant2@test.com", "staff1@test.com"})
}

// TestUserRepository_FindUsersByRole_NoResults tests when no users match role
func TestUserRepository_FindUsersByRole_NoResults(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	users, err := userRepositoryTest.FindUsersByRole("NONEXISTENT_ROLE")

	require.NoError(t, err)
	assert.Equal(t, 0, len(users))
}

// TestUserRepository_UpdateUser_Success tests successful user update
func TestUserRepository_UpdateUser_Success(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user := model.NewUser("Original Name", "1010101010", "updatetest@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	// Update user
	createdUser.Name = "Updated Name"
	createdUser.Phone = "9876543210"

	updatedUser, err := userRepositoryTest.UpdateUser(createdUser)

	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
	assert.Equal(t, "9876543210", updatedUser.Phone)
	assert.Equal(t, "updatetest@test.com", updatedUser.Email)

	defer setup.CleanupUsers([]string{"updatetest@test.com"})
}

// TestUserRepository_DeleteUser_Success tests successful user deletion
func TestUserRepository_DeleteUser_Success(t *testing.T) {
	setup.ResetTestDB()
	defer setup.ResetTestDB()

	user := model.NewUser("Delete Test", "1313131313", "deletetest@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	// Delete user
	err = userRepositoryTest.DeleteUser(createdUser)

	require.NoError(t, err)

	// Verify deletion by trying to find
	foundUser, err := userRepositoryTest.FindUserByID(createdUser.ID)
	assert.Error(t, err)
	assert.Nil(t, foundUser)

	defer setup.CleanupUsers([]string{"deletetest@test.com"})
}
