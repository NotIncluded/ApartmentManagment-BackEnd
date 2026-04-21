package model

type BillSlip struct {
	ID        string `gorm:"type:char(36);primaryKey" json:"id"`
	BillID    string `gorm:"not null" json:"bill_id"`
	RoomID    string `gorm:"not null" json:"room_id"`
	SlipURL   string `gorm:"not null" json:"slip_url"`
}
