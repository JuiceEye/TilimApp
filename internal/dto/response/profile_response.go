package response

import "time"

type GetProfileResponse struct {
	UserID           int64     `json:"user_id"`
	Username         string    `json:"username"`
	RegistrationDate time.Time `json:"registration_date"`
	Image            string    `json:"image"`
	StreakDays       int       `json:"streak_days"`
	XPPoints         int64     `json:"xp_points"`
	WordsLearned     int       `json:"words_learned"`
	LessonsDone      int       `json:"lessons_done"`
	IsSubscribed     *bool     `json:"is_subscribed,omitempty"`
}

type UserActivityResponse struct {
	Date             string `json:"date"`
	LessonsCompleted int64  `json:"lessons_completed"`
}
