package achievement

import (
	"fmt"
	"sync"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type AchievementService struct {
	achievementRepo *repository.AchievementRepository
	userRepo        *repository.UserRepository
	achievements    []Achievement
	mu              sync.RWMutex
}

func NewAchievementService(
	achievementRepo *repository.AchievementRepository,
	userRepo *repository.UserRepository,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
		achievements:    make([]Achievement, 0),
	}
}

func (s *AchievementService) RegisterAchievement(achievement Achievement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.achievements = append(s.achievements, achievement)
}

func (s *AchievementService) RegisterAchievements(achievements []Achievement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.achievements = append(s.achievements, achievements...)
}

func (s *AchievementService) Process(ctx EventContext) error {
	totalXP := 0
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, achievement := range s.achievements {
		if achievement.Trigger() != ctx.EventType {
			continue
		}

		achieved, err := achievement.Check(ctx)
		if err != nil {
			return fmt.Errorf("error checking achievement %s: %w", achievement.ID(), err)
		}

		if achieved {
			xp, err := achievement.Grant(ctx)
			if err != nil {
				return fmt.Errorf("error granting achievement %s: %w", achievement.ID(), err)
			}
			totalXP += xp
		}
	}

	err := s.userRepo.AddXP(ctx.UserID, int64(totalXP))
	if err != nil {
		return err
	}

	return nil
}

func (s *AchievementService) GetAchievements(userID int64) ([]model.Achievement, error) {
	achievements, err := s.achievementRepo.GetAchievementsWithUserStatus(userID)

	return achievements, err
}
