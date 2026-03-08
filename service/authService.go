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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(req LoginRequest, signature []byte) (string, error) {

	user, err := s.userRepo.FindUserByEmail(&req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if user.Password != req.Password {
		return "", errors.New("invalid email or password")
	}

	tokenString, err := auth.GenerateToken(signature, user.Email, user.Role)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Register(name, phone, email, password, role string) (*model.User, error) {
	existingUser, _ := s.userRepo.FindUserByEmail(&email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	newUser := model.NewUser(name, phone, email, password, role)

	return s.userRepo.CreateUser(newUser)
}
