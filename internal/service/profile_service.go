package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

type ProfileService struct {
	userRepo         *repository.UserRepository
	userProgressRepo *repository.UserProgressRepository
}

func NewProfileService(
	userRepo *repository.UserRepository,
	userProgressRepo *repository.UserProgressRepository,
) *ProfileService {
	return &ProfileService{
		userRepo:         userRepo,
		userProgressRepo: userProgressRepo,
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

	profile = &model.Profile{
		UserID:           user.ID,
		Username:         user.Username,
		Image:            user.Image,
		RegistrationDate: user.RegistrationDate,
		StreakDays:       userProgress.StreakDays,
		XPPoints:         userProgress.XPPoints,
		WordsLearned:     userProgress.WordsLearned,
		LessonsDone:      userProgress.LessonsDone,
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

func (s *ProfileService) UpdateEmail(userID int64, newEmail string) error {
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

// func (s *ProfileService) ProcessStreakTx(tx *sql.Tx, userID int64, activityDate time.Time) error {
// 	userProgress, err := s.userRepo.GetStreakTx(tx, userID)
// 	if err != nil {
// 		return err
// 	}
//
// 	if streak.LastActivity.Equal(activityDate) {
// 		return nil
// 	}
//
// 	if streak.LastActivity.Equal(activityDate.AddDate(0, 0, -1)) {
// 		streak.Current += 1
// 	} else {
// 		streak.Current = 1
// 	}
//
// 	if streak.Current > streak.Longest {
// 		streak.Longest = streak.Current
// 	}
//
// 	streak.LastActivity = activityDate
//
// 	return r.SaveStreakTx(tx, userID, streak)
// }
