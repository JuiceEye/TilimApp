package service

import (
	"fmt"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type AuthService struct {
	repository *repository.AuthRepository
}

func NewAuthService(repository *repository.AuthRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (s *AuthService) Register(user model.User) (*model.User, error, int) {
	if _, err := s.repository.GetUserByEmail(user.Email); err == nil {
		return nil, fmt.Errorf("email already taken"), 400
	}

	if _, err := s.repository.GetUserByPhone(user.PhoneNumber); err == nil {
		return nil, fmt.Errorf("phone number already taken"), 400
	}

	if _, err := s.repository.GetUserByUsername(user.Username); err == nil {
		return nil, fmt.Errorf("username already taken"), 400
	}

	createdUser, err := s.repository.CreateUser(&user)
	if err != nil {
		return nil, err, 500
	}

	return createdUser, nil, 200
}

//func (s *AuthService) Login(username, password string) (string, error) {

//}
