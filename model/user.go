package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"type:text;not null;unique"`
	Password string `json:"password" gorm:"type:text;not null"`
	Role string `json:"role" gorm:"type:text;not null;check:role IN ('ADMIN','TENANT')"`
}

func (User) TableName() string {
	return "users"
}

func NewUser(username string, password string, role string) *User {
	return &User{
		Username: username,
		Password: password,
		Role: role,
	}
}