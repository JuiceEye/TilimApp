package service

import (
	"errors"
	"net/http"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type ProfileService struct {
	userRepository         *repository.UserRepository
	userProgressRepository *repository.UserProgressRepository
}

func NewProfileService(
	userRepository *repository.UserRepository,
	userProgressRepository *repository.UserProgressRepository,
) *ProfileService {
	return &ProfileService{
		userRepository:         userRepository,
		userProgressRepository: userProgressRepository,
	}
}

func (s *ProfileService) GetProfile(userID int64) (profile *model.Profile, status int, err error) {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	userProgress, err := s.userProgressRepository.GetUserProgressByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			userProgress, err = s.userProgressRepository.CreateUserProgress(userID)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}

	profile = &model.Profile{
		UserID:           user.ID,
		Username:         user.Username,
		RegistrationDate: user.RegistrationDate,
		StreakDays:       userProgress.StreakDays,
		XPPoints:         userProgress.XPPoints,
		WordsLearned:     userProgress.WordsLearned,
		LessonsDone:      userProgress.LessonsDone,
	}

	return profile, http.StatusOK, nil
}

func (s *ProfileService) AddXPPoints(userID, xp int64) (status int, err error) {
	if _, err := s.userRepository.GetUserByID(userID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	if err := s.userProgressRepository.AddXP(userID, xp); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
