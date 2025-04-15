package response

import "time"

type ReadProfileResponse struct {
	UserID           int64     `json:"user_id"`
	Username         string    `json:"username"`
	RegistrationDate time.Time `json:"registration_date"`
	StreakDays       int       `json:"streak_days"`
	XPPoints         int       `json:"xp_points"`
	WordsLearned     int       `json:"words_learned"`
	LessonsDone      int       `json:"LessonsDone"`
}
