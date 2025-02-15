package model

import "time"

type User struct {
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Image            string    `json:"image"`
	RegistrationDate time.Time `json:"registration_date"`
}
