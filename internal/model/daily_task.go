package model

import (
	"time"
)

// DailyTask represents a task that can be assigned to users daily
type DailyTask struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	XP          int64     `json:"xp"`
	LessonID    int64     `json:"lesson_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserDailyTask represents a daily task assigned to a user
type UserDailyTask struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	DailyTaskID  int64     `json:"daily_task_id"`
	LessonID     int64     `json:"lesson_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	XP           int64     `json:"xp"`
	Completed    bool      `json:"completed"`
	AssignedDate time.Time `json:"assigned_date"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}