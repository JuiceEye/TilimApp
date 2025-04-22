package request

import "fmt"

type GetLessonRequest struct {
	LessonID int64 `json:"lesson_id"`
}

func (req *GetLessonRequest) ValidateRequest() (err error) {
	if req.LessonID <= 0 {
		return fmt.Errorf("lesson_id не может быть меньше 1")
	}

	return nil
}
