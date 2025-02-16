package dto

import (
	"fmt"
	"regexp"
)

type AuthRegistrationRequest struct {
	Username    string `json:"username"`
	Password    string `json:"Password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
	Image       string `json:"image"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

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
		missingFields = append(missingFields, "phone")
	}
	if len(missingFields) > 0 {
		err = fmt.Errorf("missing required fields: %v", missingFields)
	}
	if !phoneRegex.MatchString(req.PhoneNumber) {
		err = fmt.Errorf("invalid phone number")
	}
	if !emailRegex.MatchString(req.Email) {
		err = fmt.Errorf("invalid email address")
	}
	return err
}
