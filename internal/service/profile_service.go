package service

import (
	"errors"
	"fmt"
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

func (s *ProfileService) UpdateProfilePicture(userID int64, image string) error {
	return s.userRepository.UpdateUserImageByID(userID, image)
}

func (s *ProfileService) UpdateUsername(userID int64, newUsername string) error {
	currentUser, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if currentUser.Username == newUsername {
		return errors.New("имя пользователя должно отличаться от старого")
	}

	otherUser, err := s.userRepository.GetUserByUsername(newUsername)
	if err == nil {
		if otherUser.ID != currentUser.ID {
			return errors.New("имя пользователя уже занято")
		}
	} else if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("не удалось проверить имя пользователя: %w", err)
	}

	return s.userRepository.ChangeUserFields(userID, &model.User{Username: newUsername})
}
