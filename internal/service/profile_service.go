package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"tilimauth/internal/utils"
	"time"
)

type ProfileService struct {
	userRepo         *repository.UserRepository
	userProgressRepo *repository.UserProgressRepository
	subRepo          *repository.SubscriptionRepository
}

func NewProfileService(
	userRepo *repository.UserRepository,
	userProgressRepo *repository.UserProgressRepository,
	subRepo *repository.SubscriptionRepository,
) *ProfileService {
	return &ProfileService{
		userRepo:         userRepo,
		userProgressRepo: userProgressRepo,
		subRepo:          subRepo,
	}
}

func (s *ProfileService) GetProfile(userID int64) (profile *model.Profile, status int, err error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	userProgress, err := s.userProgressRepo.GetUserProgressByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			userProgress, err = s.userProgressRepo.CreateUserProgress(userID)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
		} else {
			return nil, http.StatusInternalServerError, err
		}
	}

	isSubscribed, err := s.subRepo.ExistsActive(userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	profile = &model.Profile{
		UserID:           user.ID,
		Username:         user.Username,
		Image:            user.Image,
		RegistrationDate: utils.ToAppTZ(user.RegistrationDate),
		StreakDays:       userProgress.StreakDays,
		XPPoints:         userProgress.XPPoints,
		WordsLearned:     userProgress.WordsLearned,
		LessonsDone:      userProgress.LessonsDone,
		IsSubscribed:     isSubscribed,
	}

	return profile, http.StatusOK, nil
}

func (s *ProfileService) UpdateProfilePicture(userID int64, image string) error {
	return s.userRepo.ChangeUserFields(userID, &model.User{Image: image})
}

func (s *ProfileService) UpdateUsername(userID int64, newUsername string) error {
	currentUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if currentUser.Username == newUsername {
		return &BadRequestError{Msg: "имя пользователя должно отличаться от старого"}
	}

	otherUser, err := s.userRepo.GetUserByUsername(newUsername)
	if err == nil {
		if otherUser.ID != currentUser.ID {
			return &BadRequestError{Msg: "имя пользователя уже занято"}
		}
	} else if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("не удалось проверить имя пользователя: %w", err)
	}

	return s.userRepo.ChangeUserFields(userID, &model.User{Username: newUsername})
}

func (s *ProfileService) UpdateEmail(userID int64, newEmail string, password string) error {
	hashedPassword, err := s.userRepo.GetUserPasswordByID(userID)
	if err != nil {
		return err
	}

	if err := utils.ComparePassword(hashedPassword, password); err != nil {
		return &BadRequestError{Msg: "неверный пароль"}
	}

	currentUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if currentUser.Email == newEmail {
		return &BadRequestError{Msg: "почта должна отличаться от старой"}
	}

	otherUser, err := s.userRepo.GetUserByEmail(newEmail)
	if err == nil {
		if otherUser.ID != currentUser.ID {
			return &BadRequestError{Msg: "почта уже используется"}
		}
	} else if !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("не удалось проверить почту: %w", err)
	}

	return s.userRepo.ChangeUserFields(userID, &model.User{Email: newEmail})
}

func (s *ProfileService) UpdatePassword(userID int64, oldPassword, newPassword string) error {
	hashedPassword, err := s.userRepo.GetUserPasswordByID(userID)
	if err != nil {
		return err
	}

	if err := utils.ComparePassword(hashedPassword, oldPassword); err != nil {
		return &BadRequestError{Msg: "неверный пароль"}
	}

	if oldPassword == newPassword {
		return &BadRequestError{Msg: "пароль должен отличаться от старого"}
	}

	newHashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("не удалось захешировать новый пароль: %w", err)
	}

	return s.userRepo.ChangeUserFields(userID, &model.User{Password: newHashedPassword})

}

func (s *ProfileService) ProcessStreakTx(tx *sql.Tx, userID int64, activityDate time.Time) error {
	userProgress, err := s.userProgressRepo.GetUserProgressByUserIDTx(tx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			userProgress, err = s.userProgressRepo.CreateUserProgress(userID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	switch {
	case userProgress.LastLessonCompletedAt != nil && userProgress.LastLessonCompletedAt.Equal(activityDate):
		return nil
	case userProgress.LastLessonCompletedAt != nil && userProgress.LastLessonCompletedAt.Equal(activityDate.AddDate(0, 0, -1)):
		userProgress.StreakDays += 1
	default:
		userProgress.StreakDays = 1
	}

	userProgress.LastLessonCompletedAt = &activityDate

	return s.userProgressRepo.SaveStreakTx(tx, userID, userProgress)
}
