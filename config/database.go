package config

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
}