package service

import "tilimauth/internal/repository"

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

//func (s *UserService) Login(username, password string) (string, error) {
//
//}
