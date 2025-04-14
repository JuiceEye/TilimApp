package model

import "time"

type Profile struct {
	UserID           int64
	Username         string
	RegistrationDate time.Time
	StreakDays       int
	XPPoints         int
	WordsLearned     int
	LessonsDone      int
}
