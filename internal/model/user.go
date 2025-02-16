package model

import "time"

type User struct {
	Id               int       `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone"`
	Image            string    `json:"image"`
	RegistrationDate time.Time `json:"registration_date"`
}
