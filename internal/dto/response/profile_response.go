package response

import "time"

type ReadProfileResponse struct {
	UserID           int64
	Username         string
	RegistrationDate time.Time
	StreakDays       int
	XPPoints         int
	WordsLearned     int
	LessonsDone      int
}
