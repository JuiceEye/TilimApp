package request

import (
	"fmt"
)

type GetProfileRequest struct {
	UserID int64 `json:"user_id"`
}

func (req *GetProfileRequest) ValidateRequest() (err error) {
	if req.UserID <= 0 {
		return fmt.Errorf("user_id не может быть меньше 1")
	}

	return nil
}
