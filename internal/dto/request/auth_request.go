package request

import (
	"errors"
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

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
var PhoneRegex = regexp.MustCompile(`^\+?[1-9]\d{9,15}$`)

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
	if !EmailRegex.MatchString(req.Email) {
		return errors.New("invalid email address")
	}
	if !PhoneRegex.MatchString(req.PhoneNumber) {
		return errors.New("invalid phone number")
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
