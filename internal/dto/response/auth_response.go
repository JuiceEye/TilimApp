package dto

type AuthRegistrationResponse struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}
