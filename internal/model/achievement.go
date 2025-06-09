package model

import (
	"time"
)

// Achievement represents an achievement that can be earned by users
type Achievement struct {
	ID          int64
	Code        string
	Name        string
	Description string
	XPReward    int
	CreatedAt   time.Time
}

// UserAchievement represents an achievement earned by a user
type UserAchievement struct {
	ID            int64
	UserID        int64
	AchievementID int64
	EarnedAt      time.Time
}