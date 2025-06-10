package model

import (
	"time"
)

type Achievement struct {
	ID          int64      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	XPReward    int        `json:"xp_reward"`
	CreatedAt   time.Time  `json:"created_at"`
	IsUnlocked  bool       `json:"is_unlocked"`
	EarnedAt    *time.Time `json:"earned_at,omitempty"`
}

type UserAchievement struct {
	ID            int64
	UserID        int64
	AchievementID int64
	EarnedAt      time.Time
}
