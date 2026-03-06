package service

import (
	"errors"

	"github.com/PunMung-66/ApartmentSys/internal/auth"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(req LoginRequest, signature []byte) (string, error) {

	user, err := s.userRepo.FindUserByUsername(&req.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if user.Password != req.Password {
		return "", errors.New("invalid username or password")
	}

	tokenString, err := auth.GenerateToken(signature, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Register(username, password string) (*model.User, error) {
	existingUser, _ := s.userRepo.FindUserByUsername(&username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	newUser := &model.User{
		Username: username,
		Password: password,
		Role:     "TENANT",
	}

	return s.userRepo.CreateUser(newUser)
}