package setup

import (
	"time"

	"github.com/PunMung-66/ApartmentSys/model"
)

// CreateTestBill helper creates a test bill for a room
func CreateTestBill(userID string, roomID string) *model.Bill {
	// For simplicity, create a contract first using the real userID
	contract, _ := CreateTestContract(userID, roomID, time.Now().Format("2006-01-02"), time.Now().AddDate(0, 6, 0).Format("2006-01-02"), "Active")
	
	// ✅ FIX: Only pass the 8 required arguments to the constructor
	bill := model.NewBill(
		contract.ID,                  // 1. Contract ID
		"test-rate-id",               // 2. Rate ID
		time.Now(),                   // 3. Record Date
		1000,                         // 4. Rent Fee
		100,                          // 5. Water Fee
		100,                          // 6. Electricity Fee
		50,                           // 7. Common Fee
		time.Now().AddDate(0, 0, 30), // 8. Due Date
	)

	// ✅ FIX: Assign the RoomID so the test database is happy
	bill.RoomID = roomID

	TestDB.Create(&bill)
	return bill
}