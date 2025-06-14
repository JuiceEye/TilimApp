package model

import (
	"time"
)

type UserProgress struct {
	UserID                int64
	StreakDays            int
	XPPoints              int64
	WordsLearned          int
	LessonsDone           int
	LastLessonCompletedAt *time.Time
	UpdatedAt             time.Time
	LastStreakResetDate   *time.Time
}
