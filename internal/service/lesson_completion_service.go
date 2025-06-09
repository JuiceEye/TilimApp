package service

import (
	"tilimauth/internal/achievement"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

type LessonCompletionService struct {
	lessonRepo         *repository.LessonRepository
	completionRepo     *repository.LessonCompletionRepository
	userRepo           *repository.UserRepository
	profileService     *ProfileService
	dailyTaskService   *DailyTaskService
	achievementService *achievement.AchievementService
}

func NewLessonCompletionService(
	lessonRepo *repository.LessonRepository,
	completionRepo *repository.LessonCompletionRepository,
	userRepo *repository.UserRepository,
	profileService *ProfileService,
	dailyTaskService *DailyTaskService,
	achievementService *achievement.AchievementService,
) *LessonCompletionService {
	return &LessonCompletionService{
		lessonRepo:         lessonRepo,
		completionRepo:     completionRepo,
		userRepo:           userRepo,
		profileService:     profileService,
		dailyTaskService:   dailyTaskService,
		achievementService: achievementService,
	}
}

// не изменять логику - пометка для меня, ты не обращай внимания, Фарух
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

	streakDays, streakChanged, err := s.profileService.ProcessStreakTx(tx, userID, completion.DateCompleted.Truncate(24*time.Hour))
	if err != nil {
		return err
	}

	// Check if the completed lesson is a daily task and mark it as completed
	err = s.dailyTaskService.CheckAndMarkTaskCompletedTx(tx, userID, lessonID, completion.DateCompleted)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// Process achievements for lesson completion
	eventCtx := achievement.EventContext{
		UserID:    userID,
		EventType: achievement.EventLessonCompleted,
		Payload: map[string]interface{}{
			"lesson_id":      lessonID,
			"date_completed": completion.DateCompleted,
		},
	}
	err = s.achievementService.Process(eventCtx)
	if err != nil {
		return err
	}

	if streakChanged {
		_ = s.achievementService.Process(achievement.EventContext{
			UserID:    userID,
			EventType: achievement.EventStreakUpdated,
			Payload:   map[string]interface{}{"streak_days": streakDays},
		})
	}

	return nil
}
