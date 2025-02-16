package repository

import (
	"database/sql"
	"errors"
	"tilimauth/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM auth.users WHERE email = $1", email)
	if err != nil {

		return nil, err
	}

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.Id == 0 {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (r *AuthRepository) GetUserByPhoneNumber(phoneNumber string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM auth.users WHERE phone_number = $1", phoneNumber)
	if err != nil {

		return nil, err
	}

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.Id == 0 {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (r *AuthRepository) GetUserByUsername(username string) (*model.User, error) {
	rows, err := r.db.Query("SELECT * FROM auth.users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}

	u := new(model.User)
	for rows.Next() {
		u, err = scanRowIntoUsers(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.Id == 0 {
		return nil, errors.New("user not found")
	}

	return u, nil
}

func (r *AuthRepository) CreateUser(user *model.User) (*model.User, error) {
	err := r.db.QueryRow(
		"INSERT INTO auth.users (username, password, email, phone_number, image, registration_date) "+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Username,
		user.Password,
		user.Email,
		user.PhoneNumber,
		user.Image,
		user.RegistrationDate,
	).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func scanRowIntoUsers(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
		&user.Image,
		&user.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
