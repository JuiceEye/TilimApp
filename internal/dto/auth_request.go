package dto

type AuthRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"Password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Image    string `json:"image"`
}
