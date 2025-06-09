package achievement

import (
	"fmt"
	"sync"
	"tilimauth/internal/repository"
)

// AchievementService handles the processing of events and checking for achievements
type AchievementService struct {
	achievementRepo *repository.AchievementRepository
	userRepo        *repository.UserRepository
	achievements    []Achievement
	mu              sync.RWMutex
}

// NewAchievementService creates a new achievement service
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

// RegisterAchievement adds an achievement to the registry
func (s *AchievementService) RegisterAchievement(achievement Achievement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.achievements = append(s.achievements, achievement)
}

// RegisterAchievements adds multiple achievements to the registry
func (s *AchievementService) RegisterAchievements(achievements []Achievement) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.achievements = append(s.achievements, achievements...)
}

// ProcessTx checks and grants achievements for the given event context within a transaction
func (s *AchievementService) Process(ctx EventContext) error {
	totalXP := 0
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, achievement := range s.achievements {
		// Skip achievements that don't match the event type
		if achievement.Trigger() != ctx.EventType {
			continue
		}

		// Check if the achievement conditions are met
		achieved, err := achievement.Check(ctx)
		if err != nil {
			return fmt.Errorf("error checking achievement %s: %w", achievement.ID(), err)
		}

		// If achieved, grant the achievement
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
