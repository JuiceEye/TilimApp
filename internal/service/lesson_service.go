package service

import (
	"fmt"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type LessonService struct {
	lessonRepo   *repository.LessonRepository
	exerciseRepo *repository.ExerciseRepository
	answerRepo   *repository.AnswerRepository
}

func NewLessonService(
	lessonRepo *repository.LessonRepository,
	exerciseRepo *repository.ExerciseRepository,
	answerRepo *repository.AnswerRepository,
) *LessonService {
	return &LessonService{
		lessonRepo:   lessonRepo,
		exerciseRepo: exerciseRepo,
		answerRepo:   answerRepo,
	}
}

func (s *LessonService) GetLessonByID(lessonID int64) (*model.Lesson, error) {
	lesson, err := s.lessonRepo.GetByID(lessonID)
	if err != nil {
		return nil, err
	}

	exercises, err := s.exerciseRepo.GetByLessonID(lesson.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sections: %w", err)
	}

	if len(exercises) == 0 {
		lesson.Exercises = []model.Exercise{}
		return lesson, nil
	}

	exerciseIDs := make([]int64, len(exercises))
	for i, exercise := range exercises {
		exerciseIDs[i] = exercise.ID
	}

	answersByExercise, err := s.answerRepo.GetByExerciseIDs(exerciseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}

	for i := range exercises {
		exercises[i].Answers = answersByExercise[exercises[i].ID]
		if exercises[i].Answers == nil {
			exercises[i].Answers = []model.Answer{}
		}
	}

	lesson.Exercises = exercises
	return lesson, nil
}
