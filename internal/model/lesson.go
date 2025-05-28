package model

type Lesson struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	XP        int64        `json:"xp"`
	Exercises []Exercise   `json:"exercises"`
	Status    LessonStatus `json:"status"`
}

type LessonStatus string

const (
	StatusLocked    LessonStatus = "Locked"
	StatusUnlocked  LessonStatus = "Unlocked"
	StatusCompleted LessonStatus = "Completed"
)
