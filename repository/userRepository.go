package repository

import (
	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	CreateUser(user *model.User) (model.User, error)
	FindUserByUsername( name * string) (model.User, error)
	UpdateUser(user *model.User) (model.User, error)
	DeleteUser(user *model.User) (error)
}

func NewUserRepository(db *gorm.DB) *UserRepository{
	return &UserRepository{db:db}
}

func (t *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	result := t.db.Create(&user)
	return user, result.Error
}

func (t *UserRepository) FindUserByUsername(name *string) (*model.User, error) {
	var user model.User
	result := t.db.Where("username = ?", name).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}