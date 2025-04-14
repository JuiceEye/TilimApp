package request

import (
	"errors"
)

type ReadProfileRequest struct {
	UserID int64 `json:"user_id"`
}

func (req *ReadProfileRequest) ValidateRequest() (err error) {
	if req.UserID <= 0 {
		return errors.New("user_id must be greater than zero")
	}

	return nil
}
