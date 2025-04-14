package response

type AuthRegistrationResponse struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}
