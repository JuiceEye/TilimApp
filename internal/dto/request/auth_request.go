package request

import (
	"fmt"
	"strings"
	"tilimauth/internal/validation"
)

type AuthRegistrationRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Image       string `json:"image"`
}

func (req *AuthRegistrationRequest) ValidateRequest() (err error) {
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
	if !validation.EmailRegex.MatchString(req.Email) {
		return fmt.Errorf("invalid email address")
	}
	if !validation.PhoneRegex.MatchString(req.PhoneNumber) {
		return fmt.Errorf("invalid phone number")
	}

	return nil
}

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req *AuthLoginRequest) ValidateRequest() (err error) {
	var missingFields []string
	if req.Username == "" {
		missingFields = append(missingFields, "username")
	}
	if req.Password == "" {
		missingFields = append(missingFields, "password")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: [%s]", strings.Join(missingFields, ", "))
	}

	return nil
}
