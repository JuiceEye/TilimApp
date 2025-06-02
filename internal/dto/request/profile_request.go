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

type UpdateUsernameRequest struct {
	Username string `json:"username"`
}

func (req *UpdateUsernameRequest) ValidateRequest() (err error) {
	if req.Username == "" {
		return fmt.Errorf("отсутствуют обязательные параметры: [username]")
	}

	return nil
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

func (req *UpdateEmailRequest) ValidateRequest() (err error) {
	if req.Email == "" {
		return fmt.Errorf("отсутствуют обязательные параметры: [email]")
	}

	return nil
}

type UpdatePasswordRequest struct {
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

func (req *UpdatePasswordRequest) ValidateRequest() (err error) {
	if req.Password == "" {
		return fmt.Errorf("отсутствуют обязательные параметры: [password, new_password]")
	}

	return nil
}
