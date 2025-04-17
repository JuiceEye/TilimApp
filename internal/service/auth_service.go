package service

import (
	"errors"
	"fmt"
	"net/http"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"tilimauth/internal/utils"
	"tilimauth/internal/validation"
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
		return nil, http.StatusBadRequest, fmt.Errorf("на эту электронную почту уже зарегистрирован аккаунт")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	if _, err = s.userRepository.GetUserByPhoneNumber(user.PhoneNumber); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("на этот номер уже зарегистрирован аккаунт")
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	if _, err = s.userRepository.GetUserByUsername(user.Username); err == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("имя пользователя уже занято")
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

func (s *AuthService) Login(usernameOrEmail, password string) (*model.User, int, error) {
	var user *model.User
	var err error

	if validation.EmailRegex.MatchString(usernameOrEmail) {
		user, err = s.userRepository.GetUserByEmail(usernameOrEmail)
	} else {
		user, err = s.userRepository.GetUserByUsername(usernameOrEmail)
	}

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, http.StatusUnauthorized, fmt.Errorf("неверные учетные данные")
		}
		return nil, http.StatusInternalServerError, err
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("неверные учетные данные")
	}

	return user, http.StatusOK, nil
}
