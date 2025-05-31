package request

import (
	"fmt"
	"strings"
	"tilimauth/internal/validation"
)

type AuthRegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Image    string `json:"image"`
}

func (req *AuthRegistrationRequest) ValidateRequest() (err error) {
	var missingFields []string
	if req.Username == "" {
		missingFields = append(missingFields, "имя пользователя")
	}
	if req.Password == "" {
		missingFields = append(missingFields, "пароль")
	}
	if req.Email == "" {
		missingFields = append(missingFields, "электронная почта")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("отсутствуют обязательные параметры: [%s]", strings.Join(missingFields, ", "))
	}
	if !validation.EmailRegex.MatchString(req.Email) {
		return fmt.Errorf("неверный формат электронной почты")
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
		missingFields = append(missingFields, "имя пользователя")
	}
	if req.Password == "" {
		missingFields = append(missingFields, "пароль")
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("отсутствуют обязательные параметры: [%s]", strings.Join(missingFields, ", "))
	}

	return nil
}
