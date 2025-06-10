package achievement

import (
	"fmt"
	"tilimauth/internal/repository"
	"time"
)

type BaseAchievement struct {
	id              string
	trigger         EventType
	achievementRepo *repository.AchievementRepository
	userRepo        *repository.UserRepository
}

func (a *BaseAchievement) ID() string {
	return a.id
}

func (a *BaseAchievement) Trigger() EventType {
	return a.trigger
}

type LessonsSingleDayAchievement struct {
	BaseAchievement
	requiredCount int
}

func NewLessonsSingleDayAchievement(
	id string,
	requiredCount int,
	achievementRepo *repository.AchievementRepository,
	userRepo *repository.UserRepository,
) *LessonsSingleDayAchievement {
	return &LessonsSingleDayAchievement{
		BaseAchievement: BaseAchievement{
			id:              id,
			trigger:         EventLessonCompleted,
			achievementRepo: achievementRepo,
			userRepo:        userRepo,
		},
		requiredCount: requiredCount,
	}
}

func (a *LessonsSingleDayAchievement) Check(ctx EventContext) (bool, error) {
	hasAchievement, err := a.achievementRepo.HasUserEarnedAchievementByCode(ctx.UserID, a.id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has achievement: %w", err)
	}

	if hasAchievement {
		return false, nil
	}

	var date time.Time
	if dateVal, ok := ctx.Payload["date_completed"]; ok {
		if dateTime, ok := dateVal.(time.Time); ok {
			date = dateTime
		} else {
			date = time.Now().UTC()
		}
	} else {
		date = time.Now().UTC()
	}

	count, err := a.achievementRepo.GetLessonCompletionsCountForDay(ctx.UserID, date)
	if err != nil {
		return false, fmt.Errorf("failed to get lesson completions count: %w", err)
	}

	return count >= a.requiredCount, nil
}

func (a *LessonsSingleDayAchievement) Grant(ctx EventContext) (awardableXP int, err error) {
	achievement, err := a.achievementRepo.GetAchievementByCode(a.id)
	if err != nil {
		return 0, fmt.Errorf("failed to get achievement: %w", err)
	}

	tx, err := a.userRepo.DB().Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = a.achievementRepo.GrantAchievementTx(tx, ctx.UserID, achievement.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to grant achievement: %w", err)
	}

	return achievement.XPReward, tx.Commit()
}

type LessonsTotalAchievement struct {
	BaseAchievement
	requiredCount int
}

func NewLessonsTotalAchievement(
	id string,
	requiredCount int,
	achievementRepo *repository.AchievementRepository,
	userRepo *repository.UserRepository,
) *LessonsTotalAchievement {
	return &LessonsTotalAchievement{
		BaseAchievement: BaseAchievement{
			id:              id,
			trigger:         EventLessonCompleted,
			achievementRepo: achievementRepo,
			userRepo:        userRepo,
		},
		requiredCount: requiredCount,
	}
}

func (a *LessonsTotalAchievement) Check(ctx EventContext) (bool, error) {
	hasAchievement, err := a.achievementRepo.HasUserEarnedAchievementByCode(ctx.UserID, a.id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has achievement: %w", err)
	}

	if hasAchievement {
		return false, nil
	}

	count, err := a.achievementRepo.GetTotalLessonCompletionsCount(ctx.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to get total lesson completions count: %w", err)
	}

	return count >= a.requiredCount, nil
}

func (a *LessonsTotalAchievement) Grant(ctx EventContext) (awardableXP int, err error) {
	achievement, err := a.achievementRepo.GetAchievementByCode(a.id)
	if err != nil {
		return 0, fmt.Errorf("failed to get achievement: %w", err)
	}

	tx, err := a.userRepo.DB().Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = a.achievementRepo.GrantAchievementTx(tx, ctx.UserID, achievement.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to grant achievement: %w", err)
	}

	return achievement.XPReward, tx.Commit()
}

type LessonStreakAchievement struct {
	BaseAchievement
	requiredDays int
}

func NewLessonStreakAchievement(
	id string,
	requiredDays int,
	achievementRepo *repository.AchievementRepository,
	userRepo *repository.UserRepository,
) *LessonStreakAchievement {
	return &LessonStreakAchievement{
		BaseAchievement: BaseAchievement{
			id:              id,
			trigger:         EventStreakUpdated,
			achievementRepo: achievementRepo,
			userRepo:        userRepo,
		},
		requiredDays: requiredDays,
	}
}

func (a *LessonStreakAchievement) Check(ctx EventContext) (bool, error) {
	hasAchievement, err := a.achievementRepo.HasUserEarnedAchievementByCode(ctx.UserID, a.id)
	if err != nil {
		return false, fmt.Errorf("failed to check if user has achievement: %w", err)
	}

	if hasAchievement {
		return false, nil
	}

	var streakDays int
	if streakVal, ok := ctx.Payload["streak_days"]; ok {
		if streak, ok := streakVal.(int); ok {
			streakDays = streak
		}
	}

	return streakDays >= a.requiredDays, nil
}

func (a *LessonStreakAchievement) Grant(ctx EventContext) (awardableXP int, err error) {
	achievement, err := a.achievementRepo.GetAchievementByCode(a.id)
	if err != nil {
		return 0, fmt.Errorf("failed to get achievement: %w", err)
	}

	tx, err := a.userRepo.DB().Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = a.achievementRepo.GrantAchievementTx(tx, ctx.UserID, achievement.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to grant achievement: %w", err)
	}

	return achievement.XPReward, tx.Commit()
}

func CreateDefaultAchievements(
	achievementService *AchievementService,
	achievementRepo *repository.AchievementRepository,
	userRepo *repository.UserRepository,
) {
	achievements := []Achievement{
		NewLessonStreakAchievement("LESSON_STREAK_3", 3, achievementRepo, userRepo),
		NewLessonStreakAchievement("LESSON_STREAK_7", 7, achievementRepo, userRepo),
		NewLessonStreakAchievement("LESSON_STREAK_30", 30, achievementRepo, userRepo),

		NewLessonsSingleDayAchievement("LESSONS_SINGLE_DAY_5", 5, achievementRepo, userRepo),

		NewLessonsTotalAchievement("LESSONS_TOTAL_10", 10, achievementRepo, userRepo),
		NewLessonsTotalAchievement("LESSONS_TOTAL_50", 50, achievementRepo, userRepo),
		NewLessonsTotalAchievement("LESSONS_TOTAL_100", 100, achievementRepo, userRepo),
	}

	achievementService.RegisterAchievements(achievements)
}
