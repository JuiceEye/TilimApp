package service

import (
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

type LessonCompletionService struct {
	lessonRepo     *repository.LessonRepository
	completionRepo *repository.LessonCompletionRepository
	userRepo       *repository.UserRepository
	profileService *ProfileService
}

func NewLessonCompletionService(
	lessonRepo *repository.LessonRepository,
	completionRepo *repository.LessonCompletionRepository,
	userRepo *repository.UserRepository,
	profileService *ProfileService,
) *LessonCompletionService {
	return &LessonCompletionService{
		lessonRepo:     lessonRepo,
		completionRepo: completionRepo,
		userRepo:       userRepo,
		profileService: profileService,
	}
}

func (s *LessonCompletionService) CompleteLesson(completion *model.LessonCompletion) error {
	userID := completion.UserID
	lessonID := completion.LessonID
	lesson, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return err
	}

	exists, err := s.completionRepo.Exists(userID, lessonID)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	tx, err := s.userRepo.DB().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.completionRepo.RegisterTx(tx, completion)
	if err != nil {
		return err
	}

	err = s.userRepo.IncrementStatsTx(tx, userID, lesson.XP)
	if err != nil {
		return err
	}

	err = s.profileService.ProcessStreakTx(tx, userID, completion.DateCompleted.Truncate(24*time.Hour))
	if err != nil {
		return err
	}

	return tx.Commit()
}
