package response

import (
	"time"
	"tilimauth/internal/model"
)

// GetDailyTaskResponse represents a daily task response
type GetDailyTaskResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	XP          int64      `json:"xp"`
	LessonID    int64      `json:"lesson_id"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// ToDailyTaskResponseList converts a list of UserDailyTask models to a list of GetDailyTaskResponse DTOs
func ToDailyTaskResponseList(tasks []model.UserDailyTask) []GetDailyTaskResponse {
	var responses []GetDailyTaskResponse
	for _, task := range tasks {
		responses = append(responses, GetDailyTaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			XP:          task.XP,
			LessonID:    task.LessonID,
			Completed:   task.Completed,
			CompletedAt: task.CompletedAt,
		})
	}
	return responses
}