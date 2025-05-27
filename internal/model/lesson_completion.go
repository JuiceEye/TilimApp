package model

import "time"

type LessonCompletion struct {
	UserID        int64
	LessonID      int64
	DateCompleted time.Time
}
