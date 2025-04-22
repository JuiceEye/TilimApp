package request

import "fmt"

type LessonCompletedRequest struct {
	UserID   int64 `json:"user_id"`
	XPPoints int64 `json:"xp_points"`
}

func (r *LessonCompletedRequest) ValidateRequest() error {
	if r.UserID <= 0 {
		return fmt.Errorf("user_id не может быть меньше 1")
	}
	if r.XPPoints < 0 {
		return fmt.Errorf("xp_points не могут быть отрицательным числом")
	}
	return nil
}
