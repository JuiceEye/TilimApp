package service

import (
	"errors"
	"fmt"
	"net/http"
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

func (s *AuthService) Register(user model.User) (createdUser *model.User, status int, err error) {
	if _, err := s.repository.GetUserByEmail(user.Email); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("email already taken")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	if _, err = s.repository.GetUserByPhoneNumber(user.PhoneNumber); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("phone number already taken")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	if _, err = s.repository.GetUserByUsername(user.Username); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("username already taken")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	createdUser, err = s.repository.CreateUser(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return createdUser, http.StatusOK, nil
}

//func (s *AuthService) Login(username, password string) (string, error) {

//}
