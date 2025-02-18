package dto

import (
	"fmt"
	"regexp"
	"strings"
)

type AuthRegistrationRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Image       string `json:"image"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{9,15}$`)

func (req *AuthRegistrationRequest) Validate() error {
	var err error
	var missingFields []string
	if req.Username == "" {
		missingFields = append(missingFields, "username")
	}
	if req.Password == "" {
		missingFields = append(missingFields, "password")
	}
	if req.Email == "" {
		missingFields = append(missingFields, "email")
	}
	if req.PhoneNumber == "" {
		missingFields = append(missingFields, "phone_number")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: [%s]", strings.Join(missingFields, ", "))
	}
	if !emailRegex.MatchString(req.Email) {
		err = fmt.Errorf("invalid email address")
	}
	if !phoneRegex.MatchString(req.PhoneNumber) {
		err = fmt.Errorf("invalid phone number")
	}
	return err
}

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
