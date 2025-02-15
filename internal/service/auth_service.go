package service

import "tilimauth/internal/repository"

type AuthService struct {
	repository *repository.AuthRepository
}

func NewAuthService(repository *repository.AuthRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

//func (s *AuthService) Login(username, password string) (string, error) {

//}
