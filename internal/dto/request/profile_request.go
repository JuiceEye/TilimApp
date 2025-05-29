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

type UpdateProfilePictureRequest struct {
	Image string `json:"image"`
}

func (req *UpdateProfilePictureRequest) ValidateRequest() (err error) {
	if req.Image == "" {
		return fmt.Errorf("отсутствуют обязательные параметры: [image]")
	}

	return nil
}
