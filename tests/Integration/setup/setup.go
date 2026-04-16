package setup

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/PunMung-66/ApartmentSys/config"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	TestDBName = "apartment_test"
	JWTSecret  = "test_jwt_secret_key"
)

// Type aliases for easier access
type (
	User     = model.User
	Room     = model.Room
	Contract = model.Contract
)

var (
	TestDB       *gorm.DB
	UserRepo     *repository.UserRepository
	RoomRepo     *repository.RoomRepository
	ContractRepo *repository.ContractRepository
	AuthService  *service.AuthService
	UserService  *service.UserService
	RoomService  *service.RoomService
	Env          string
)

func init() {
	flag.StringVar(&Env, "e", "local", "Environment: local | dev")
	// Auto-initialize database when setup package is imported by any subpackage
	// This ensures tests work even when run with: go test ./tests/Integration/controller
	if TestDB == nil {
		InitTestDatabase(getEnv())
	}
}

// getEnv returns the test environment, defaulting to "local"
func getEnv() string {
	if Env != "" {
		return Env
	}
	return "local"
}

// InitTestDatabase initializes the test database connection
func InitTestDatabase(environment string) {
	var host, port, user, password string

	switch environment {
	case "dev":
		host = "3.80.58.141"
		port = "5432"
		user = "postgres"
		password = "devpassword"

	case "local":
		host = "localhost"
		port = "8080"
		user = "postgres"
		password = "mysecretpassword"

	default:
		panic("Invalid environment: " + environment)
	}

	os.Setenv("DB_HOST", host)
	os.Setenv("DB_USER", user)
	os.Setenv("DB_PASSWORD", password)
	os.Setenv("DB_PORT", port)

	db, err := config.ConnectTestDatabase(TestDBName)
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Room{})
	db.AutoMigrate(&model.Contract{})

	TestDB = db
	UserRepo = repository.NewUserRepository(TestDB)
	RoomRepo = repository.NewRoomRepository(TestDB)
	ContractRepo = repository.NewContractRepository(TestDB)
	AuthService = service.NewAuthService(UserRepo)
	UserService = service.NewUserService(UserRepo)
	RoomService = service.NewRoomService(RoomRepo, ContractRepo)
}

// ResetTestDB clears all test data
func ResetTestDB() {
	if TestDB != nil {
		TestDB.Exec("TRUNCATE TABLE contracts CASCADE")
		TestDB.Exec("TRUNCATE TABLE rooms CASCADE")
		TestDB.Exec("TRUNCATE TABLE users CASCADE")
	}
}

// TeardownTestDB closes the database connection
func TeardownTestDB() {
	if TestDB != nil {
		TestDB.Exec("TRUNCATE TABLE contracts CASCADE")
		TestDB.Exec("TRUNCATE TABLE rooms CASCADE")
		TestDB.Exec("TRUNCATE TABLE users CASCADE")
		sqlDB, err := TestDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// CleanupUsers removes specific users by email
func CleanupUsers(emails []string) {
	if TestDB == nil {
		panic("TestDB is nil - database not initialized")
	}
	for _, email := range emails {
		TestDB.Unscoped().Where("email = ?", email).Delete(&model.User{})
	}
}

// CleanupRooms removes specific rooms by ID
func CleanupRooms(roomIDs []string) {
	if TestDB == nil {
		panic("TestDB is nil - database not initialized")
	}
	for _, roomID := range roomIDs {
		TestDB.Unscoped().Where("id = ?", roomID).Delete(&model.Room{})
	}
}

// CreateTestRoom helper creates a test room
func CreateTestRoom(roomNumber string, level int, status string) *model.Room {
	if TestDB == nil {
		panic("TestDB is nil - database not initialized. Make sure testmain_test.go exists in Integration package root")
	}
	room := model.NewRoom(roomNumber, level, status)
	result := TestDB.Create(&room)
	if result.Error != nil {
		panic("Failed to create test room: " + result.Error.Error())
	}
	return room
}

// CreateTestContract helper creates a test contract
func CreateTestContract(userID, roomID string, startDate, endDate string, status string) (*model.Contract, error) {
	if TestDB == nil {
		panic("TestDB is nil - database not initialized. Make sure testmain_test.go exists in Integration package root")
	}

	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, errors.New("invalid start date format")
		}
	} else {
		start = time.Now()
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, errors.New("invalid end date format")
		}
	} else {
		end = time.Now().AddDate(0, 6, 0)
	}

	contract := &model.Contract{
		UserID:    userID,
		RoomID:    roomID,
		StartDate: start,
		EndDate:   end,
		Status:    status,
	}

	result := TestDB.Create(&contract)
	if result.Error != nil {
		return nil, result.Error
	}
	return contract, nil
}
