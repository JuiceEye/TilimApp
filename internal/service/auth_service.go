package service

import (
	"errors"
	"fmt"
	"net/http"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type AuthService struct {
	userRepository         *repository.UserRepository
	userProgressRepository *repository.UserProgressRepository
}

func NewAuthService(
	userRepository *repository.UserRepository,
	userProgressRepository *repository.UserProgressRepository,
) *AuthService {
	return &AuthService{
		userRepository:         userRepository,
		userProgressRepository: userProgressRepository,
	}
}

func (s *AuthService) Register(user model.User) (createdUser *model.User, status int, err error) {
	if _, err := s.userRepository.GetUserByEmail(user.Email); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("email already taken")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	if _, err = s.userRepository.GetUserByPhoneNumber(user.PhoneNumber); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("phone number already taken")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	if _, err = s.userRepository.GetUserByUsername(user.Username); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("username already taken")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	// user.Password = Bcrypt(user.Password)
	createdUser, err = s.userRepository.CreateUser(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	_, err = s.userProgressRepository.CreateUserProgress(createdUser.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return createdUser, http.StatusOK, nil
}

// func (s *AuthService) Login(username, password string) (string, error) {

// }
