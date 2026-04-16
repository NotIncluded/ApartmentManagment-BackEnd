package setup

import (
       "time"
       "github.com/PunMung-66/ApartmentSys/model"
)

// CreateTestBill helper creates a test bill for a room
func CreateTestBill(roomID string) *model.Bill {
       // For simplicity, create a contract first
       contract, _ := CreateTestContract("test-user-id", roomID, time.Now().Format("2006-01-02"), time.Now().AddDate(0, 6, 0).Format("2006-01-02"), "Active")
       bill := model.NewBill(contract.ID, "test-rate-id", time.Now(), 1000, 100, 100, 50, 1250, "Unpaid", time.Now().AddDate(0, 0, 30), time.Now())
       TestDB.Create(&bill)
       return bill
}