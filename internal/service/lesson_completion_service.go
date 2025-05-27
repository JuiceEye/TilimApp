package service

import (
	"errors"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type LessonCompletionService struct {
	lessonRepo     *repository.LessonRepository
	completionRepo *repository.LessonCompletionRepository
}

func NewLessonCompletionService(
	lessonRepo *repository.LessonRepository,
	completionRepo *repository.LessonCompletionRepository,
) *LessonCompletionService {
	return &LessonCompletionService{
		lessonRepo:     lessonRepo,
		completionRepo: completionRepo,
	}
}

func (s *LessonCompletionService) CompleteLesson(lessonCompletion *model.LessonCompletion) error {
	_, err := s.lessonRepo.GetByID(lessonCompletion.LessonID)
	if err != nil {
		return err
	}

	if isExists, err := s.completionRepo.Exists(lessonCompletion.UserID, lessonCompletion.LessonID); err != nil {
		return err
	} else {
		if isExists {
			return errors.New("lesson is already completed")
		}
	}

	if err = s.completionRepo.Register(lessonCompletion); err != nil {
		return err
	}

	return nil
}
