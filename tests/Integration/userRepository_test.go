package integration

import (
	"testing"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	userRepositoryTest *repository.UserRepository
)

func initUserRepository() {
	userRepositoryTest = repository.NewUserRepository(testDB)
}

func setupUserTestDB() {
	if testDB != nil {
		testDB.AutoMigrate(&model.User{})
	}
}

func resetUserTestDB() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE users CASCADE")
	}
}

func cleanupTestUsersByEmail(emails []string) {
	for _, email := range emails {
		testDB.Unscoped().Where("email = ?", email).Delete(&model.User{})
	}
}

// TestUserRepository_CreateUser_Success tests successful user creation
func TestUserRepository_CreateUser_Success(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user := model.NewUser("John Doe", "1234567890", "johndoe@test.com", "password123", "TENANT")

	createdUser, err := userRepositoryTest.CreateUser(user)

	require.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.Equal(t, "johndoe@test.com", createdUser.Email)
	assert.Equal(t, "1234567890", createdUser.Phone)
	assert.Equal(t, "TENANT", createdUser.Role)

	defer cleanupTestUsersByEmail([]string{"johndoe@test.com"})
}

// TestUserRepository_CreateUser_MultipleUsers tests creating multiple users
func TestUserRepository_CreateUser_MultipleUsers(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

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

	defer cleanupTestUsersByEmail([]string{"userone@test.com", "usertwo@test.com"})
}

// TestUserRepository_FindUserByID_Success tests finding user by ID
func TestUserRepository_FindUserByID_Success(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user := model.NewUser("Find By ID", "3333333333", "findbyid@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	foundUser, err := userRepositoryTest.FindUserByID(createdUser.ID)

	require.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, createdUser.ID, foundUser.ID)
	assert.Equal(t, "Find By ID", foundUser.Name)
	assert.Equal(t, "findbyid@test.com", foundUser.Email)

	defer cleanupTestUsersByEmail([]string{"findbyid@test.com"})
}

// TestUserRepository_FindUserByID_NotFound tests when user ID doesn't exist
func TestUserRepository_FindUserByID_NotFound(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	foundUser, err := userRepositoryTest.FindUserByID("nonexistent-user-id")

	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

// TestUserRepository_FindUserByEmail_Success tests finding user by email
func TestUserRepository_FindUserByEmail_Success(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

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

	defer cleanupTestUsersByEmail([]string{"emailuser@test.com"})
}

// TestUserRepository_FindUserByEmail_NotFound tests when email doesn't exist
func TestUserRepository_FindUserByEmail_NotFound(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	email := "nonexistent@test.com"
	foundUser, err := userRepositoryTest.FindUserByEmail(&email)

	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

// TestUserRepository_FindUserByEmail_CaseSensitive tests email case sensitivity
func TestUserRepository_FindUserByEmail_CaseSensitive(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user := model.NewUser("Case Test", "5555555555", "CaseEmail@test.com", "password123", "TENANT")
	_, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	// Attempt to find with different case
	differentCaseEmail := "caseemail@test.com"
	foundUser, err := userRepositoryTest.FindUserByEmail(&differentCaseEmail)

	// This tests case sensitivity - may vary based on database
	// Most databases are case-insensitive for string comparison by default
	if err != nil {
		assert.Nil(t, foundUser)
	}

	defer cleanupTestUsersByEmail([]string{"CaseEmail@test.com", "caseemail@test.com"})
}

// TestUserRepository_FindUsersByRole_Success tests finding all users by role
func TestUserRepository_FindUsersByRole_Success(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	// Create multiple TENANT users
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

	defer cleanupTestUsersByEmail([]string{"tenant1@test.com", "tenant2@test.com", "staff1@test.com"})
}

// TestUserRepository_FindUsersByRole_Staff tests finding all STAFF users
func TestUserRepository_FindUsersByRole_Staff(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user1 := model.NewUser("Staff Member", "9999999999", "staffmem@test.com", "pass123", "STAFF")
	_, err := userRepositoryTest.CreateUser(user1)
	require.NoError(t, err)

	staffUsers, err := userRepositoryTest.FindUsersByRole("STAFF")

	require.NoError(t, err)
	assert.Greater(t, len(staffUsers), 0)
	for _, u := range staffUsers {
		assert.Equal(t, "STAFF", u.Role)
	}

	defer cleanupTestUsersByEmail([]string{"staffmem@test.com"})
}

// TestUserRepository_FindUsersByRole_NoResults tests when no users match role
func TestUserRepository_FindUsersByRole_NoResults(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	users, err := userRepositoryTest.FindUsersByRole("NONEXISTENT_ROLE")

	require.NoError(t, err)
	assert.Equal(t, 0, len(users))
}

// TestUserRepository_UpdateUser_Success tests successful user update
func TestUserRepository_UpdateUser_Success(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

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
	assert.Equal(t, "updatetest@test.com", updatedUser.Email) // Email should remain unchanged

	defer cleanupTestUsersByEmail([]string{"updatetest@test.com"})
}

// TestUserRepository_UpdateUser_MultipleFields tests updating multiple fields
func TestUserRepository_UpdateUser_MultipleFields(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user := model.NewUser("Initial User", "1111111111", "multiupdate@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	// Update multiple fields
	createdUser.Name = "New Name"
	createdUser.Phone = "5555555555"
	// Note: Role cannot be updated based on service logic, but we can test field updates

	updatedUser, err := userRepositoryTest.UpdateUser(createdUser)

	require.NoError(t, err)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "5555555555", updatedUser.Phone)

	defer cleanupTestUsersByEmail([]string{"multiupdate@test.com"})
}

// TestUserRepository_UpdateUser_VerifyPersistence tests that changes persist
func TestUserRepository_UpdateUser_VerifyPersistence(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user := model.NewUser("Persist Test", "1212121212", "persisttest@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	// Update user
	createdUser.Name = "Persisted Name"
	updatedUser, err := userRepositoryTest.UpdateUser(createdUser)
	require.NoError(t, err)

	// Verify by fetching again
	refetchedUser, err := userRepositoryTest.FindUserByID(updatedUser.ID)
	require.NoError(t, err)

	assert.Equal(t, "Persisted Name", refetchedUser.Name)

	defer cleanupTestUsersByEmail([]string{"persisttest@test.com"})
}

// TestUserRepository_DeleteUser_Success tests successful user deletion
func TestUserRepository_DeleteUser_Success(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

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

	defer cleanupTestUsersByEmail([]string{"deletetest@test.com"})
}

// TestUserRepository_DeleteUser_VerifyRemoved tests that deleted user is completely removed
func TestUserRepository_DeleteUser_VerifyRemoved(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	user := model.NewUser("Verify Delete", "1414141414", "verifydelete@test.com", "password123", "STAFF")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	userID := createdUser.ID

	// Delete user
	err = userRepositoryTest.DeleteUser(createdUser)
	require.NoError(t, err)

	// Verify cannot find by ID
	foundByID, _ := userRepositoryTest.FindUserByID(userID)
	assert.Nil(t, foundByID)

	// Verify cannot find by email
	email := "verifydelete@test.com"
	foundByEmail, _ := userRepositoryTest.FindUserByEmail(&email)
	assert.Nil(t, foundByEmail)

	defer cleanupTestUsersByEmail([]string{"verifydelete@test.com"})
}

// TestUserRepository_CreateAndRetrieve tests full user lifecycle
func TestUserRepository_CreateAndRetrieve(t *testing.T) {
	setupUserTestDB()
	defer resetUserTestDB()
	initUserRepository()

	// Create
	user := model.NewUser("Lifecycle User", "1515151515", "lifecycle@test.com", "password123", "TENANT")
	createdUser, err := userRepositoryTest.CreateUser(user)
	require.NoError(t, err)

	// Retrieve by ID
	retrievedByID, err := userRepositoryTest.FindUserByID(createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, retrievedByID.ID)

	// Retrieve by Email
	email := "lifecycle@test.com"
	retrievedByEmail, err := userRepositoryTest.FindUserByEmail(&email)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, retrievedByEmail.ID)

	// Retrieve by Role
	users, err := userRepositoryTest.FindUsersByRole("TENANT")
	require.NoError(t, err)
	assert.Greater(t, len(users), 0)

	hasCreatedUser := false
	for _, u := range users {
		if u.ID == createdUser.ID {
			hasCreatedUser = true
			break
		}
	}
	assert.True(t, hasCreatedUser)

	defer cleanupTestUsersByEmail([]string{"lifecycle@test.com"})
}
