package request

import (
	"errors"
)

type ReadUserRequest struct {
	UserId int64 `json:"user_id"`
}

func (req *ReadUserRequest) ValidateRequest() (err error) {
	if req.UserId <= 0 {
		return errors.New("user_id must be greater than zero")
	}

	return nil
}
