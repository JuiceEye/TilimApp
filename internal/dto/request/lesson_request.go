package request

import "fmt"

type GetLessonRequest struct {
	LessonID int64 `json:"lesson_id"`
}

type CompleteLessonRequest struct {
	LessonID int64 `json:"lesson_id"`
	UserID   int64 `json:"user_id"`
}

func (req *GetLessonRequest) ValidateRequest() (err error) {
	if req.LessonID <= 0 {
		return fmt.Errorf("lesson_id не может быть меньше 1")
	}

	return nil
}

func (req *CompleteLessonRequest) ValidateRequest() (err error) {
	if req.LessonID <= 0 {
		return fmt.Errorf("lesson_id не может быть меньше 1")
	}

	return nil
}
