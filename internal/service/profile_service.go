package service

import (
	"errors"
	"net/http"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type ProfileService struct {
	repository *repository.ProfileRepository
}

func NewProfileService(repository *repository.ProfileRepository) *ProfileService {
	return &ProfileService{
		repository: repository,
	}
}

func (s *ProfileService) GetProfile(userID int64) (profile *model.Profile, status int, err error) {
	profile, err = s.repository.GetProfileById(userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, http.StatusNotFound, err
	} else if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return profile, http.StatusOK, nil
}
