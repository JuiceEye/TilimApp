package service

import (
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type LessonCompletionService struct {
	lessonRepo     *repository.LessonRepository
	completionRepo *repository.LessonCompletionRepository
	userRepo       *repository.UserRepository
}

func NewLessonCompletionService(
	lessonRepo *repository.LessonRepository,
	completionRepo *repository.LessonCompletionRepository,
	userRepo *repository.UserRepository,
) *LessonCompletionService {
	return &LessonCompletionService{
		lessonRepo:     lessonRepo,
		completionRepo: completionRepo,
		userRepo:       userRepo,
	}
}

func (s *LessonCompletionService) CompleteLesson(lessonCompletion *model.LessonCompletion) error {
	userID := lessonCompletion.UserID
	lessonID := lessonCompletion.LessonID
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

	err = s.completionRepo.RegisterTx(tx, lessonCompletion)
	if err != nil {
		return err
	}

	err = s.userRepo.IncrementStatsTx(tx, userID, lesson.XP)
	if err != nil {
		return err
	}

	return tx.Commit()

}
