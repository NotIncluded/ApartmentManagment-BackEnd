package repository

import (
	"github.com/PunMung-66/ApartmentSys/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	CreateUser(user *model.User) (*model.User, error)
	FindUserByID(id string) (*model.User, error)
	FindUserByEmail(email *string) (*model.User, error)
	FindUsersByRole(role string) ([]model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(user *model.User) error
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (t *UserRepository) CreateUser(user *model.User) (*model.User, error) {
	result := t.db.Create(&user)
	return user, result.Error
}

func (t *UserRepository) FindUserByID(id string) (*model.User, error) {
	var user model.User
	result := t.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (t *UserRepository) FindUserByEmail(email *string) (*model.User, error) {
	var user model.User
	result := t.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (t *UserRepository) FindUsersByRole(role string) ([]model.User, error) {
	var users []model.User
	result := t.db.Where("role = ?", role).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (t *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	result := t.db.Save(&user)
	return user, result.Error
}

func (t *UserRepository) DeleteUser(user *model.User) error {
	result := t.db.Delete(&user)
	return result.Error
}
