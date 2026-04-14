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

func (s *UserService) GetUserByID(userID string) (*model.User, error) {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUsersByRole(role string) ([]model.User, error) {
	users, err := s.userRepo.FindUsersByRole(role)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {
	if user.Name == "" || user.Email == "" {
		return nil, errors.New("incomplete request body")
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(userID string) error {
	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return err
	}

	return s.userRepo.DeleteUser(user)
}
