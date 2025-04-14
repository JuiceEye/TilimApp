package model

import "time"

type User struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone_number"`
	Image            string    `json:"image"`
	RegistrationDate time.Time `json:"registration_date"`
}
