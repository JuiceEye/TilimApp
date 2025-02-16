package dto

type AuthRegistrationResponse struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}
