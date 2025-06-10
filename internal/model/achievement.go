package model

import (
	"time"
)

type Achievement struct {
	ID          int64
	Code        string
	Name        string
	Description string
	XPReward    int
	CreatedAt   time.Time
}

type UserAchievement struct {
	ID            int64
	UserID        int64
	AchievementID int64
	EarnedAt      time.Time
}
