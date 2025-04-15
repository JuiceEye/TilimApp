package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
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

func (s *AuthService) Login(usernameOrEmail, password string) (*model.User, int, error) {
	var user *model.User
	var err error

	if isEmail(usernameOrEmail) {
		user, err = s.userRepository.GetUserByEmail(usernameOrEmail)
	} else {
		user, err = s.userRepository.GetUserByUsername(usernameOrEmail)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, http.StatusUnauthorized, fmt.Errorf("invalid credentials")
		}
		return nil, http.StatusInternalServerError, err
	}

	//if err := utils.ComparePassword(user.HashedPassword, password); err != nil {
	//	return nil, http.StatusUnauthorized, fmt.Errorf("invalid credentials")
	//}

	return user, http.StatusOK, nil
}

func isEmail(input string) bool {
	return strings.Contains(input, "@") && strings.Contains(input, ".") //чекаем есть ли собачка и точка - главные идентификаторы имейла

}
