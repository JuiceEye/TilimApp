package response

type AccessRefreshTokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type AuthRegistrationResponse struct {
	UserID int64                   `json:"user_id"`
	Tokens *AccessRefreshTokenPair `json:"tokens"`
}

type AuthLoginResponse struct {
	UserID int64                   `json:"user_id"`
	Tokens *AccessRefreshTokenPair `json:"tokens"`
}

type TokenRefreshResponse struct {
	UserID int64                   `json:"user_id"`
	Tokens *AccessRefreshTokenPair `json:"tokens"`
}
