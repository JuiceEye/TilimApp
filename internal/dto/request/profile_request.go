package request

import (
	"fmt"
)

type ReadProfileRequest struct {
	UserID int64 `json:"user_id"`
}

func (req *ReadProfileRequest) ValidateRequest() (err error) {
	if req.UserID <= 0 {
		return fmt.Errorf("user_id не может быть меньше 1")
	}

	return nil
}
