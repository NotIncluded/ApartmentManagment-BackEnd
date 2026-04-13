package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"user_id"`
	Name      string         `json:"name" gorm:"not null"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"password" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;check:role IN ('ADMIN','TENANT')"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

func NewUser(name, phone, email, password, role string) *User {
	return &User{
		Name:     name,
		Phone:    phone,
		Email:    email,
		Password: password,
		Role:     role,
	}
}
