package service

import (
	"errors"

	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
func (s *UserService) CreateUser(user *model.User) (*model.User, error) {

	if user.Name == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		return nil, errors.New("incomplete request body")
	}

	userResponse, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (s *UserService) DeleteUser(userID string) error {
	user := &model.User{}
	user.ID = userID

	err := s.userRepo.DeleteUser(user)
	if err != nil {
		return err
	}

	return nil
}
