package service

import (
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

func (s AuthService) Register(user model.User) (model.User, error, int) {
	return user, nil, 0
}

//func (s *AuthService) Login(username, password string) (string, error) {

//}
